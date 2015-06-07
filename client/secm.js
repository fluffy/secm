/*jslint  browser: true, devel: true,  vars: true, todo: true, white: true, bitwise: true */
/*global  $,WebSocket */
/* jshint strict: true, jquery: true */

/* TODO 
clean up 
error check all JSON parse
Make it all Promises
 */

var Fluffy;
Fluffy = Fluffy || {}; // setup namespace

Fluffy.SecM = (function () {
    "use strict";

    function setup() {
        if (!crypto.subtle) {
            alert("No WebCrypto subtle support");
            return;
        }
    }

    function arrayToString(a) {
        return String.fromCharCode.apply(null, new Uint8Array(a));
    }

    function arrayToHexString(a) {
        var lRes = new Uint8Array(a);
        var ret = "";
        for (var i in lRes) {
            var v = lRes[i];
            ret += (v < 16 ? "0" : "") + v.toString(16);
        }
        return ret;
    }

    function hexStringToArray(s) {
        console.assert(s.length % 2 === 0, "string must be of even length to covert");
        var i;
        var l = s.length / 2;
        var ret = new ArrayBuffer(l);
        var view = new Uint8Array(ret);
        for (i = 0; i < l; i++) {
            view[i] = parseInt(s.substring(i * 2, (i + 1) * 2), 16);
        }
        return ret;
    }

    function stringToArray(s) {
        var i, len;
        var ret = new ArrayBuffer(s.length);
        var view = new Uint8Array(ret);
        for (i = 0, len = s.length; i < len; i++) {
            view[i] = s.charCodeAt(i);
        }
        return ret;
    }

    function genKey(f) {
        crypto.subtle.generateKey({
                name: "AES-GCM",
                length: 128
            },
            true, ["encrypt", "decrypt"]
        ).then(
            function (key) {
                crypto.subtle.exportKey(
                    "jwk",
                    key
                ).then(function (expKey) {
                    var stringKey = JSON.stringify(expKey);
                    f(stringKey);
                }).
                catch (function (err) {
                    console.log("problem exporting key: " + err);
                });

            }).
        catch (function (err) {
            console.log("problem generating key: " + err);
        });
    }

    function checkSign(encString, jwkObj, f) {
        console.assert($.type(encString) === "string", "encrypt takes string");
        console.assert($.type(jwkObj) === "object", "encrypt takes string");

        var encObj = JSON.parse(encString); // todo - move out and error check 
        var aad = stringToArray(encObj.authData);
        var tag = hexStringToArray(encObj.tag);
        var iv = hexStringToArray(encObj.iv);

        crypto.subtle.importKey(
            "jwk",
            jwkObj, {
                name: "AES-GCM",
                length: 128
            },
            true, ["encrypt", "decrypt"]
        ).then(function (key) {
            crypto.subtle.decrypt({
                    name: "AES-GCM",
                    additionalData: aad, // optional 
                    tagLength: 128, // 128,104,32,64,96,112,120  // optional (128 or missing , len=32 ) (32, len=20)
                    iv: iv // required 
                },
                key,
                tag
            ).then(function (dResR) {
                f(encObj.authData);
            }).
            catch (function (err) {
                console.log("problem checking sig: " + err);
            });

        }).
        catch (function (err) {
            console.log("problem importing sig key: " + err);
        });
    }

    function decrypt(encString, jwkObj, f) {
        console.assert($.type(encString) === "string", "encrypt takes string");
        console.assert($.type(jwkObj) === "object", "encrypt takes string");

        var encObj = JSON.parse(encString); // todo - move out and error check 
        var aad = new Uint8Array([]);
        var cipherText = hexStringToArray(encObj.ct);
        var iv = hexStringToArray(encObj.iv);

        crypto.subtle.importKey(
            "jwk",
            jwkObj, {
                name: "AES-GCM",
                length: 128
            },
            true, ["encrypt", "decrypt"]
        ).then(function (key) {
            crypto.subtle.decrypt({
                    name: "AES-GCM",
                    additionalData: aad, // optional 
                    tagLength: 32, // 128,104,32,64,96,112,120  // optional (128 or missing , len=32 ) (32, len=20)
                    iv: iv // required 
                },
                key,
                cipherText
            ).then(function (result) {
                var resString = arrayToString(result);
                f(resString);

            }).
            catch (function (err) {
                console.log("problem decrypting : " + err);
            });

        }).
        catch (function (err) {
            console.log("problem importing decrypt key: " + err);
        });
    }

    function encrypt(dataString, jwkObj, f) {
        console.assert($.type(dataString) === "string", "encrypt takes string");
        console.assert($.type(jwkObj) === "object", "encrypt takes string");

        var data = stringToArray(dataString);
        var iv = crypto.getRandomValues(new Uint8Array(12));
        var aad = new Uint8Array([]);

        crypto.subtle.importKey(
            "jwk",
            jwkObj, {
                name: "AES-GCM",
                length: 128
            },
            true, ["encrypt", "decrypt"]
        ).then(function (key) {
            crypto.subtle.encrypt({
                    name: "AES-GCM",
                    additionalData: aad, // optional 
                    tagLength: 32, // 32,64,96,104,112,120,128
                    iv: iv // required 
                },
                key,
                data
            ).then(function (encText) {
                // encText is the encrypted data with the tag concatinated to it 

                var result = {
                    iv: arrayToHexString(iv),
                    ct: arrayToHexString(encText)
                };

                f(JSON.stringify(result));

            }).
            catch (function (err) {
                console.log("problem encrypting: " + err);
            });

        }).
        catch (function (err) {
            console.log("problem importing encrypt key: " + err);
        });
    }

    function sign(dataString, jwkObj, f) {
        console.assert($.type(dataString) === "string", "encrypt takes string");
        console.assert($.type(jwkObj) === "object", "encrypt takes string");

        var aad = stringToArray(dataString);
        var iv = crypto.getRandomValues(new Uint8Array(12));
        var data = new Uint8Array([]);

        crypto.subtle.importKey(
            "jwk",
            jwkObj, {
                name: "AES-GCM",
                length: 128
            },
            true, ["encrypt", "decrypt"]
        ).then(function (key) {
            crypto.subtle.encrypt({
                    name: "AES-GCM",
                    additionalData: aad, // optional 
                    tagLength: 128, // 32,64,96,104,112,120,128
                    iv: iv // required 
                },
                key,
                data
            ).then(function (res) {
                // res is the encrypted data with the tag concatinated to it 
                var result = {
                    iv: arrayToHexString(iv),
                    authData: dataString,
                    tag: arrayToHexString(res)
                };

                f(JSON.stringify(result));

            }).
            catch (function (err) {
                console.log("problem encrypting for signature: " + err);
            });

        }).
        catch (function (err) {
            console.log("problem importing key for signature: " + err);
        });

    }

    var publicExport = {
        setup: setup,
        genKey: genKey,
        encrypt: encrypt,
        sign: sign,
        checkSign: checkSign,
        decrypt: decrypt
    };

    return publicExport;
}());

$(document).ready(function () {
    "use strict";
    Fluffy.SecM.setup();

    $("#keyID").val("758614435099350414"); // TODO remove 
    $("#msgIn").val("abc"); // TODO remove
    $("#seqNum").val("7"); // TODO remove

    $("#genBut").click(function () {
        Fluffy.SecM.genKey(function (key) {
            $("#uKeyIn").val(key);
        });
        Fluffy.SecM.genKey(function (key) {
            $("#cKeyIn").val(key);
        });
    });

    $("#storeKeyBut").click(function () {
        $.post($("#ksUrl").val() + "v1/key", {
                keyVal: $("#uKeyIn").val(),
                iKeyVal: $("#cKeyIn").val()
            },
            function (data, status) {
                var d = jQuery.parseJSON(data);
                $("#keyID").val(d.keyID);
            });
    });

    $("#fetchUKeyBut").click(function () {
        $.get($("#ksUrl").val() + "v1/key/" + $("#keyID").val(),
            function (data, status) {
                $("#uKeyOut").val(data);
            });
    });

    $("#fetchCKeyBut").click(function () {
        $.get($("#ksUrl").val() + "v1/iKey/" + $("#keyID").val(),
            function (data, status) {
                $("#cKeyOut").val(data);
            });
    });

    $("#encMsgBut").click(function () {
        var jwkObj = JSON.parse($("#uKeyOut").val());
        var data = $("#msgIn").val();

        Fluffy.SecM.encrypt(data, jwkObj, function (s) {
            $("#msgEnc").val(s);
        });
    });

    $("#signMsgBut").click(function () {
        var jwkObj = JSON.parse($("#cKeyOut").val());
        var data = $("#msgEnc").val();

        Fluffy.SecM.sign(data, jwkObj, function (s) {
            $("#msgSign").val(s);
        });
    });

    $("#postMsgBut").click(function () {
        $.ajax({
            type: "POST",
            url: $("#msUrl").val() + "v1/ch/" + $("#keyID").val(),
            data: $("#msgSign").val(),
            success: function (data, status) {
                var d = jQuery.parseJSON(data);
                $("#seqNum").val(d.seqNum);
            },
            contentType: false
        });
    });

    $("#fetchMsgBut").click(function () {
        $.get($("#msUrl").val() + "v1/msg/" + $("#keyID").val() + "-" + $("#seqNum").val(),
            function (data, status) {
                $("#msgOut").val(data);
            });
    });

    $("#unSignMsgBut").click(function () {
        var jwkObj = {};
        try {
            jwkObj = JSON.parse($("#cKeyOut").val());
        } catch (err) {
            console.log("Error parsing JSON=" + $("#cKeyOut").val());
        }

        var data = $("#msgOut").val();

        Fluffy.SecM.checkSign(data, jwkObj, function (s) {
            $("#msgEncOut").val(s);
        });
    });

    $("#decryptMsgBut").click(function () {
        var jwkObj = {};
        try {
            jwkObj = JSON.parse($("#uKeyOut").val());
        } catch (err) {
            console.log("Error parsing JSON=" + $("#uKeyOut").val());
        }

        var data = $("#msgEncOut").val();

        Fluffy.SecM.decrypt(data, jwkObj, function (s) {
            $("#msgOutDecrypt").val(s);
        });
    });
});
