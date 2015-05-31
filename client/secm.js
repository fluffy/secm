
$(document).ready(function(){
    $("#createID").click(function(){
        $.post(   $("#ksUrlID").val() + "v1/key",
        {
            keyVal: $("#inValID").val(),
            iKeyVal: $("#inVal2ID").val()
        },
        function(data,status){
           var d = jQuery.parseJSON( data );
           $("#keyID").val( d.keyID );
        });
  });
  $("#getKeyID").click(function(){
        $.get( $("#ksUrlID").val() + "v1/key/" + $("#keyID").val() ,
        function(data,status){
           $("#outKeyValID").text( data );
        });
    });
  $("#getIKeyID").click(function(){
        $.get( $("#ksUrlID").val() + "v1/iKey/" + $("#keyID").val() ,
        function(data,status){
           $("#outIKeyValID").text( data );
        });
  });

  $("#postMsgID").click(function(){
        $.ajax({ type: "POST",
                 url : $("#msUrlID").val() + "v1/ch/" + $("#keyID").val(),
                 data: $("#msgInVal").val(),
                 success: function(data,status){
                     var d = jQuery.parseJSON( data );
                     $("#seqNumID").val( d.seqNum );
                 },
                 contentType: false
        });
  });

  $("#getMsgID").click(function(){
        $.get( $("#msUrlID").val() + "v1/msg/" + $("#keyID").val() + "-" + $("#seqNumID").val() ,
        function(data,status){
           $("#outValID").text( data );
        });
  });
});
