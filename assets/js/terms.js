$(document).ready(function() {

  /**
   * Update single translation and
   */
  var fnUpdateTranslation = function() {
      var $termId = $(this).data("term-id");
      var data = {
        value: $(this).val(),
        id: $termId
      };

      $.post( "/api/terms/" + $termId, data)
        .done(function( data ) {
          $('#savel-label-' + $termId)
              .transition('scale')
              .transition('scale');
        });
  };

  $(".clickable-term").click(function() {
    var $termId = $(this).data("id");

    $("#translations-panel").empty();
    $("#translations-panel").load( "/api/terms/" + $termId, function() {
      $(".saved-label").hide();
      $(".row").removeClass("active");
      $("#row-" + $termId).addClass("active");

      // on blur, save changes
      $(".editable").blur(fnUpdateTranslation);
    });

  });

});
