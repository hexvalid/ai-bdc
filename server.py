#!/bin/python3
import cv2
from torchvision.transforms.functional import to_tensor
from flask import Flask
from flask import request
from collections import OrderedDict
import numpy as np
import torch
import torch.nn as nn
from waitress import serve

listenaddr = '0.0.0.0'
model_file = 'model.bin'
characters = 'abcdefghijklmnopqrstuvwxyz'

class CNNCTC(nn.Module):
    def __init__(self, n_classes):
        super(CNNCTC, self).__init__()
        channels = [32, 64, 128, 256, 256]
        layers = [2, 2, 2, 2, 2]
        kernels = [3, 3, 3, 3, 3]
        pools = [2, 2, 2, 2, (2, 1)]
        modules = OrderedDict()

        def cba(name, in_channels, out_channels, kernel_size):
            modules[f'conv{name}'] = nn.Conv2d(in_channels, out_channels, kernel_size,
                                               padding=(1, 1) if kernel_size == 3 else 0)
            modules[f'bn{name}'] = nn.BatchNorm2d(out_channels)
            modules[f'relu{name}'] = nn.ReLU(inplace=True)

        last_channel = 1
        for block, (n_channel, n_layer, n_kernel, k_pool) in enumerate(zip(channels, layers, kernels, pools)):
            for layer in range(1, n_layer + 1):
                cba(f'{block + 1}{layer}', last_channel, n_channel, n_kernel)
                last_channel = n_channel
            modules[f'pool{block + 1}'] = nn.MaxPool2d(k_pool)
        modules[f'dropout'] = nn.Dropout(0.25, inplace=True)

        self.cnn = nn.Sequential(modules)
        self.lstm = nn.LSTM(input_size=self.infer_features(), hidden_size=128, num_layers=2, bidirectional=True)
        self.fc = nn.Linear(in_features=256, out_features=n_classes)

    def infer_features(self):
        x = torch.zeros((1, 1, 64, 200))
        x = self.cnn(x)
        x = x.reshape(x.shape[0], -1, x.shape[-1])
        return x.shape[1]

    def forward(self, x):
        x = self.cnn(x)
        x = x.reshape(x.shape[0], -1, x.shape[-1])
        x = x.permute(2, 0, 1)
        x, _ = self.lstm(x)
        x = self.fc(x)
        return x


def load_weights(target, source_state):
    new_dict = OrderedDict()
    for k, v in target.state_dict().items():
        if k in source_state and v.size() == source_state[k].size():
            new_dict[k] = source_state[k]
        else:
            new_dict[k] = v
    target.load_state_dict(new_dict)


def load_model(device):
    model = CNNCTC(n_classes=len(characters)).to(device)
    load_weights(model, torch.load(model_file, map_location='cpu'))
    return model


def decode(sequence):
    a = ''.join([characters[x] for x in sequence])
    s = ''.join([x for j, x in enumerate(a[:-1]) if x != characters[0] and x != a[j + 1]])
    if len(s) == 0:
        return ''
    if a[-1] != characters[0] and s[-1] != a[-1]:
        s += a[-1]
    return s


def decode_target(sequence):
    return ''.join([characters[x] for x in sequence]).replace(' ', '')
    
def recognize(img):
    lower = np.array([0, 0, 0])
    upper = np.array([100, 100, 100])
    img = cv2.inRange(img, lower, upper)
    element = cv2.getStructuringElement(cv2.MORPH_RECT, (2, 2))
    img = 255 - cv2.dilate(img, element, iterations=2)
    img = to_tensor(img).unsqueeze(0)
    pred = model(img).squeeze().argmax(-1)
    pred = decode(pred)
    return pred


app = Flask(__name__)

@app.route('/solve/', methods=['POST'])
def index():
    image = request.files['image'].read()
    npimage = np.frombuffer(image, np.uint8)
    cvimage = cv2.imdecode(npimage, cv2.IMREAD_UNCHANGED)
    answer = recognize(cvimage)
    return answer

if __name__ == '__main__':
    print('Loading model: '+model_file)
    model = load_model('cpu')
    model.eval()
    serve(app, host=listenaddr, port=4200)
