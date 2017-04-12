$(document).ready(function() {

  $(".clickable-term").click(function() {
    var $termId = $(this).data("id");
    $.get( "/api/terms/" + $termId, function() {
      alert( "success" );
    })
    .fail(function() {
        $("#translations-panel").text("Error! Can't load details");
    });
  });

});
