$(document).ready(function(){
    $('button').click(function(){
        var tempInput = document.createElement("input");
        tempInput.style = "position: absolute; left: -99999px; top: -99999px";
        tempInput.value = $.trim($("font").text());
        document.body.appendChild(tempInput);
        tempInput.select();
        document.execCommand("copy");
        document.body.removeChild(tempInput);
    });

    

});