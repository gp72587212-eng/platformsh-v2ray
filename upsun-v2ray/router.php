<?php
// ============================================================
//  V2Ray xHTTP Reverse Proxy — Upsun / Platform.sh
//  Remplace TON_IP_VPS par l'IP de ton serveur V2Ray
// ============================================================
define('VPS_HOST', 'TON_IP_VPS');
define('VPS_PORT', 443);
define('VPS_SCHEME', 'https');

@ini_set('output_buffering', 'Off');
@ini_set('zlib.output_compression', 0);
@ini_set('implicit_flush', 1);
ob_implicit_flush(true);
if (ob_get_level()) ob_end_flush();

$method = $_SERVER['REQUEST_METHOD'] ?? 'GET';
$uri    = $_SERVER['REQUEST_URI']    ?? '/';
$body   = file_get_contents('php://input');
$target = VPS_SCHEME . '://' . VPS_HOST . ':' . VPS_PORT . $uri;

$forwardHeaders = [];
foreach (getallheaders() as $k => $v) {
    $lower = strtolower($k);
    if (in_array($lower, [
        'host','connection','transfer-encoding','te','trailer',
        'upgrade','proxy-connection','keep-alive','proxy-authorization'
    ])) continue;
    $forwardHeaders[] = "$k: $v";
}
$forwardHeaders[] = 'Host: ' . VPS_HOST;
$forwardHeaders[] = 'Accept-Encoding: identity';

$ch = curl_init($target);
curl_setopt_array($ch, [
    CURLOPT_CUSTOMREQUEST  => $method,
    CURLOPT_HTTPHEADER     => $forwardHeaders,
    CURLOPT_POSTFIELDS     => $body,
    CURLOPT_WRITEFUNCTION  => function($ch, $data) {
        echo $data;
        flush();
        return strlen($data);
    },
    CURLOPT_HEADERFUNCTION => function($ch, $header) {
        $h     = trim($header);
        $lower = strtolower($h);
        if ($h === '') return strlen($header);
        foreach (['transfer-encoding','connection','keep-alive',
                  'proxy-connection','upgrade','te','trailer'] as $s) {
            if (str_starts_with($lower, $s . ':')) return strlen($header);
        }
        if (str_starts_with($lower, 'http/')) {
            $parts = explode(' ', $h, 3);
            http_response_code((int)($parts[1] ?? 200));
            return strlen($header);
        }
        header($h, false);
        return strlen($header);
    },
    CURLOPT_SSL_VERIFYPEER => false,
    CURLOPT_SSL_VERIFYHOST => 0,
    CURLOPT_CONNECTTIMEOUT => 10,
    CURLOPT_TIMEOUT        => 0,
    CURLOPT_HTTP_VERSION   => CURL_HTTP_VERSION_1_1,
    CURLOPT_BINARYTRANSFER => true,
    CURLOPT_BUFFERSIZE     => 16384,
    CURLOPT_FOLLOWLOCATION => false,
]);

if (in_array($method, ['GET', 'HEAD'])) {
    curl_setopt($ch, CURLOPT_POSTFIELDS, null);
    if ($method === 'HEAD') curl_setopt($ch, CURLOPT_NOBODY, true);
}

$ok  = curl_exec($ch);
$err = curl_error($ch);
curl_close($ch);

if (!$ok && $err) {
    http_response_code(502);
    echo "Proxy error: $err";
}
