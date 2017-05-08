$(document).ready(function() {

  $projectId = $('#projectId').val();

  /**
   * Activate selected term
   */
  var fnClickTermCallback = function() {
    var $termId = $(this).data("id");

    $("#translations-panel").empty();
    $("#translations-panel").load( "/api/terms/" + $termId, function() {
      $(".saved-label").hide();
      $(".item").removeClass("active");
      $("#row-" + $termId).addClass("active");

      // on blur, save changes
      $(".editable").blur(fnUpdateTranslation);
    });

 };

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
          $('#savel-label-' + $termId + "-" + $lang)
              .transition('scale')
              .transition('scale');
        });
  };

  /**
   * Add new language request
   */
  var fnAddNewLanguage = function(selectedCountryCode) {
    var data = {
      projectId: $projectId,
      countryCode: selectedCountryCode,
      isDefault: false
    };

    $.post( "/api/project/add/language", data)
      .done(function( data ) {
        location.reload();
      });
  };

  /**
   * Request to post new term to server
   */
  var fnAddNewTerm = function() {
    var $termKey = $('#new-term-key').val();
    var $termDescr = $('#new-term-description').val();

    var data = {
      projectId: $projectId,
      termKey: $termKey,
      termDescr: $termDescr
    };

    $.post( "/api/project/add/term", data)
      .done(function( data ) {

        // term was added, lets add it to the common list and activate it
        $('#new-term-key').val('');
        $('#new-term-description').val('');

        $('#table-terms > tbody:last-child').append(
          '<tr id="row-' + data.term.ID + '" class="row">' +
            '<td id="row-td-' + data.term.ID + '" class="clickable-term" data-id="' + data.term.ID + '">' + data.term.Code + '</td>' +
          '</tr>');
        $('#row-td-' + data.term.ID).click(fnClickTermCallback);
      });
  };

  /**
   * On a term selected, we load its translations on the right
   */
  $(".clickable-term").click(fnClickTermCallback);

 /**
  * Show the dialog add new language
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

 /**
  * Show the dialog "add new term"
  */
 $('#add-term-button').click(function() {

    $('#add-new-term-form').form({
      fields: {
        'key': 'empty'
      },
      on: 'blur'
    });

   $('#modal-add-term')
       .modal('setting', {
           onApprove: function () {
             var isValid = $('#add-new-term-form').form('is valid');
             if (isValid) {
                fnAddNewTerm();
             }
             else
                return false;
           }
        })
       .modal('show');

 });

 /**
  * Show the dialog "Upload new file" to import new translation
  */
 $('#import-new-file-button').click(function() {

   $('#modal-import-language')
       .modal('setting', {
           onApprove: function () {
             var isValid = $('#add-new-term-form').form('is valid');
             if (isValid) {
                fnAddNewTerm();
             }
             else
                return false;
           }
        })
       .modal('show');
   $('#country-upload-dropdown').dropdown();
 });

});
