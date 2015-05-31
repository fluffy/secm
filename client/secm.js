/*jslint  browser: true, devel: true,  vars: true, todo: true, white: true, bitwise: true */
/*global  $,WebSocket */
/* jshint strict: true, jquery: true */

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
        
    var publicExport = {
        setup: setup,
        genKey: genKey
    };
    
    return publicExport;
}());

$(document).ready(function(){
    "use strict";
    Fluffy.SecM.setup();

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
        $("#msgEnc").val(  "enc-" + $("#msgIn").val() )
    });
    
    $("#signMsgBut").click(function(){
         $("#msgSign").val(  "sign-" + $("#msgEnc").val() )
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
         $("#msgOutDecrypt").val(  "decrypt-" + $("#msgOut").val() )
    });
                         
});
