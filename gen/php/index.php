<?php session_start(); ?>

<?php require("botdetect.php");
?>

<?php
$DefaultCaptcha = new Captcha('DefaultCaptcha');
$DefaultCaptcha->Locale = 'en-US';
$DefaultCaptcha->UserInputID = 'DefaultCaptcha';
$DefaultCaptcha->ImageWidth = 240;
$DefaultCaptcha->ImageHeight = 90;
$DefaultCaptcha->CodeLength = 6;
$DefaultCaptcha->CodeStyle = CodeStyle::Alpha;
$DefaultCaptcha->ImageFormat = ImageFormat::Jpeg;
$DefaultCaptcha->HelpLinkEnabled = false;
$DefaultCaptcha->SoundEnabled = false;
$DefaultCaptcha->ReloadEnabled = false;
echo " S:";
echo $DefaultCaptcha->get_CaptchaSoundUrl();
echo ":E ";
echo $DefaultCaptcha->Html();
echo "\n<br>END";
?>
