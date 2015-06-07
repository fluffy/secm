/*jslint  browser: true, devel: true,  vars: true, todo: true, white: true, bitwise: true */
/*global  $,WebSocket */
/* jshint strict: true, jquery: true */

/* TODO 
add signing
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

        // TODO cleanup 
        //var t = arrayToHexString( stringToArray( "abc" ) );
        //console.log( "abc as hex string: " + t );
        // console.log( "hex string back to string: " + arrayToString( hexStringToArray( t ) ) );
        
    }

    function arrayToString(a) {
        return String.fromCharCode.apply( null, new Uint8Array(a) );
    }

    function arrayToHexString(a) {
        var lRes = new Uint8Array( a );
        var s = "";
        for ( var i in lRes)  {
            var v = lRes[i];
            s += (v < 16 ? "0" : "") + v.toString(16);
        }
        return s;
    }

    function hexStringToArray(s) {
        console.assert( s.length%2 === 0 , "string must be of even length to covert" );
        var i;
        var l = s.length/2;
        var ret = new ArrayBuffer( l );
        var view = new Uint8Array( ret );
        for (  i=0; i<l; i++) {
            view[i] = parseInt( s.substring( i*2, (i+1)*2 ), 16)
        }
        return ret;
    }
    
    function stringToArray(s) {
        var i,len;
        var ret = new ArrayBuffer( s.length );
        var view = new Uint8Array( ret );
        for (  i=0, len=s.length; i<len; i++) {
            view[i] = s.charCodeAt(i);
        }
        return ret;
    }

    function genKey( f ) {
        crypto.subtle.generateKey(
            { name : "AES-GCM", length:128 },
            true,
            ["encrypt","decrypt"]
        ).then(
            function(key) {
                var myKey = key;
                console.log( "gen key: " + key );
                
                crypto.subtle.exportKey(
                    "jwk",
                    myKey
                ).then(function(ekey) {
                    var sKey = JSON.stringify(ekey)
                    console.log( "Exported Key: " + sKey );
                    f( sKey );
                }).catch( function(err) {
                    console.log( "problem exporting key: " + err );
                });
                
            }).catch( function(err) {
                console.log( "problem generating key: " + err );
            });
    }

    function decrypt( encString, jwkObj, f ) {
        console.assert( $.type( encString ) === "string", "encrypt takes string");
        console.assert( $.type( jwkObj ) === "object", "encrypt takes string");

        var encObj = JSON.parse( encString ); // todo - move out and error check 

        var iKey=undefined;
        //var n = new Uint8Array([1,2,3,4,5,6,7,8,9,10,11,12]); // 96 bit IV 
        //var n = crypto.subtle.getRandomValues(new Uint8Array(12));
        var aad = new Uint8Array( [] );


        var lRes = hexStringToArray( encObj.ct );
        var n = hexStringToArray( encObj.iv );
        
        crypto.subtle.importKey(
            "jwk",
            jwkObj, 
            { name : "AES-GCM", length:128 },
            true,
            ["encrypt","decrypt"]
        ).then(function(key) {
            console.log( "Import Key: " + key );
            iKey = key;
            
            crypto.subtle.exportKey(
                "jwk",
                iKey
            ).then(function(ekey) {
                console.log( "the imported key looks like: " + JSON.stringify(ekey) );
            }).catch( function(err) {
                console.log( "problem exporting key: " + err );
            });
            
            crypto.subtle.decrypt(
                {
                    name : "AES-GCM" ,
                    additionalData: aad, // optional 
		            tagLength: 32, // 128,104,32,64,96,112,120  // optional (128 or missing , len=32 ) (32, len=20)
                    iv: n // required 
                },
                iKey,
                lRes
            ).then(function(dResR) {
                    
                console.log( "the decrypted stuff length: " + dResR.byteLength );
                var dRes = new Uint8Array( dResR );
                var s = "";
                for ( var i in dRes)  {
                    var v = dRes[i];
                    s += (v < 10 ? "0" : "") + v.toString(16);
                }
                console.log( "the decrypted stuff: " + s );
                
                var resString = arrayToString( dResR );
                console.log( "decrypted: " + resString );
                
                f( resString );        
                
            }).catch( function(err) {
                console.log( "problem decrypting : " + err );
            });
            
        }).catch( function(err) {
            console.log( "problem importing key: " + err );
            console.log( "jwk=" + JSON.stringify(jwkObj) );
        });
    }

    
    function encrypt( dataString, jwkObj, f ) {
        console.assert( $.type( dataString ) === "string", "encrypt takes string");
        console.assert( $.type( jwkObj ) === "object", "encrypt takes string");

        var iKey=undefined;
        var data = stringToArray( dataString );
        //var n = new Uint8Array([1,2,3,4,5,6,7,8,9,10,11,12]); // 96 bit IV
        var n = crypto.getRandomValues(new Uint8Array(12));
        var aad = new Uint8Array( [] );
        
        crypto.subtle.importKey(
            "jwk",
            jwkObj, 
            { name : "AES-GCM", length:128 },
            true,
            ["encrypt","decrypt"]
        ).then(function(key) {
            
            console.log( "Import Key: " + key );
            iKey = key;
            
            crypto.subtle.exportKey(
                "jwk",
                iKey
            ).then(function(ekey) {
                console.log( "the imported key looks like: " + JSON.stringify(ekey) );
            }).catch( function(err) {
                console.log( "problem exporting key: " + err );
            });
            
            
            crypto.subtle.encrypt(
                {
                    name : "AES-GCM" ,
                    additionalData: aad, // optional 
		            tagLength: 32, // 32,64,96,104,112,120,128
                    iv: n // required 
                },
                iKey,
                data
            ).then(function(res) {
                // res is the encrypted data with the tag concatinated to it 
             
                var result = {
                    iv: arrayToHexString(n),
                    ct: arrayToHexString(res)
                };
                
                f( JSON.stringify(result) );
                              
            }).catch( function(err) {
                console.log( "problem encrypting : " + err );
            });
            
        }).catch( function(err) {
            console.log( "problem importing key: " + err );
        });
    }

    function sign( dataString, jwkObj, f ) {
        console.assert( $.type( dataString ) === "string", "encrypt takes string");
        console.assert( $.type( jwkObj ) === "object", "encrypt takes string");

        var iKey=undefined;
        var aad = stringToArray( dataString );
        //var n = new Uint8Array([1,2,3,4,5,6,7,8,9,10,11,12]); // 96 bit IV
        var n = crypto.getRandomValues(new Uint8Array(12));
        var data = new Uint8Array( [] );
        
        crypto.subtle.importKey(
            "jwk",
            jwkObj, 
            { name : "AES-GCM", length:128 },
            true,
            ["encrypt","decrypt"]
        ).then(function(key) {
            
            console.log( "Import Key: " + key );
            iKey = key;
            
            crypto.subtle.exportKey(
                "jwk",
                iKey
            ).then(function(ekey) {
                console.log( "the imported key looks like: " + JSON.stringify(ekey) );
            }).catch( function(err) {
                console.log( "problem exporting key: " + err );
            });
            
            
            crypto.subtle.encrypt(
                {
                    name : "AES-GCM" ,
                    additionalData: aad, // optional 
		            tagLength: 128, // 32,64,96,104,112,120,128
                    iv: n // required 
                },
                iKey,
                data
            ).then(function(res) {
                // res is the encrypted data with the tag concatinated to it 
             
                var result = {
                    iv: arrayToHexString(n),
                    authData: dataString,
                    tag: arrayToHexString(res)
                };
                
                f( JSON.stringify(result) );
                              
            }).catch( function(err) {
                console.log( "problem encrypting sig: " + err );
            });
            
        }).catch( function(err) {
            console.log( "problem importing key for sig: " + err );
        });

    }


    var publicExport = {
        setup: setup,
        genKey: genKey,
        encrypt: encrypt,
        sign: sign,
        decrypt: decrypt
    };
    
    return publicExport;
}());

$(document).ready(function(){
    "use strict";
    Fluffy.SecM.setup();

    $("#keyID").val( "758614435099350414" ); // TODO remove 
    $("#msgIn").val( "abc" ); // TODO remove
    $("#seqNum").val( "2" );// TODO remove
    
    $("#genBut").click(function(){
        Fluffy.SecM.genKey( function(key) { $("#uKeyIn").val( key ) } );
        Fluffy.SecM.genKey( function(key) { $("#cKeyIn").val( key ) } );
    });

    $("#storeKeyBut").click(function(){
        $.post(   $("#ksUrl").val() + "v1/key",
        {
            keyVal: $("#uKeyIn").val(),
            iKeyVal: $("#cKeyIn").val()
        },
        function(data,status){
           var d = jQuery.parseJSON( data );
           $("#keyID").val( d.keyID );
        });
    });
    
    $("#fetchUKeyBut").click(function(){
        $.get( $("#ksUrl").val() + "v1/key/" + $("#keyID").val() ,
               function(data,status){
                   $("#uKeyOut").val( data );
               });
    });
    
    $("#fetchCKeyBut").click(function(){
        $.get( $("#ksUrl").val() + "v1/iKey/" + $("#keyID").val() ,
               function(data,status){
                   $("#cKeyOut").val( data );
               });
    });

    
    
    $("#encMsgBut").click(function(){
        var jwkObj = JSON.parse( $("#uKeyOut").val() );
        var data = $("#msgIn").val();
        
        Fluffy.SecM.encrypt( data, jwkObj, function( s ) { $("#msgEnc").val( s ); } );
    });
    
    $("#signMsgBut").click(function(){
        //$("#msgSign").val(  $("#msgEnc").val() ) // TODO remove
        var jwkObj = JSON.parse( $("#cKeyOut").val() );
        var data = $("#msgEnc").val();
        
        Fluffy.SecM.sign( data, jwkObj, function( s ) { $("#msgSign").val( s ); } );
    });
    
    $("#postMsgBut").click(function(){
          $.ajax({ type: "POST",
                 url : $("#msUrl").val() + "v1/ch/" + $("#keyID").val(),
                 data: $("#msgSign").val(),
                 success: function(data,status){
                     var d = jQuery.parseJSON( data );
                     $("#seqNum").val( d.seqNum );
                 },
                 contentType: false
        });
    });
    
    $("#fetchMsgBut").click(function(){
        $.get( $("#msUrl").val() + "v1/msg/" + $("#keyID").val() + "-" + $("#seqNum").val() ,
               function(data,status){
                   $("#msgOut").val( data );
               });
    });
    
    $("#decryptMsgBut").click(function(){
       // $("#msgOutDecrypt").val(  "decrypt-" + $("#msgOut").val() ) // TODO remove
        var jwkObj = {}
        try {
            jwkObj = JSON.parse( $("#uKeyOut").val() );
        } catch (err) {
            console.log( "Error parsing JSON=" + $("#uKeyOut").val() );
        }
        
        var data = $("#msgOut").val();
        
        Fluffy.SecM.decrypt( data, jwkObj, function( s ) { $("#msgOutDecrypt").val( s ); } );
    });
                         
});
