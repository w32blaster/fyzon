$(document).ready(function() {

  /**
   * Update single translation and
   */
  var fnUpdateTranslation = function() {
      var $lang = $(this).data("country-code");
      var $termId = $("#termId").val();
      var data = {
        value: $(this).val(),
        id: $termId
      };

      // this is existing term, let's update it
      $.post( "/api/terms/" + $termId + "/" + $lang, data)
        .done(function( data ) {
          $('#savel-label-' + $termId)
              .transition('scale')
              .transition('scale');
        });
  };

  /**
   * Add new language request
   */
  var fnAddNewLanguage = function(selectedCountryCode) {
    $projectId = $('#projectId').val();
    var data = {
      projectId: $projectId,
      countryCode: selectedCountryCode
    };

    $.post( "/api/project/add/language", data)
      .done(function( data ) {
        location.reload();
      });
  };

  /**
   * On a term selected, we load its translations on the right
   */
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

 /**
  * Show the dialog
  */
 $("#add-language-button").click(function() {
   $('#modal-add-lang')
      .modal('setting', {
        onApprove: function () {
          var val = $('#country-dropdown').dropdown('get value');
          fnAddNewLanguage(val);
        }
       })
      .modal('show');

    $('#country-dropdown').dropdown();
 });

});
