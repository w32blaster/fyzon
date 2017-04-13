$(document).ready(function() {

  $(".clickable-term").click(function() {
    var $termId = $(this).data("id");

    $("#translations-panel").empty();
    $("#translations-panel").load( "/api/terms/" + $termId, function() {
      $(".row").removeClass("active");
      $("#row-" + $termId).addClass("active");
      console.log("added");
    });

  });

});
