$(document).ready(function() {

  $projectId = $('#projectId').val();

  $('.icon-default-language').popup();

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
      var $isDefault = JSON.parse($(this).data("is-default")); // use JSON.parse to get bool from string

      var $termId = $("#termId").val();
      var data = {
        value: $(this).val(),
        id: $termId
      };

      // this is existing term, let's update it
      $.post( "/api/terms/" + $termId + "/" + $lang, data)
        .done(function( data ) {

          // show animated "Saved" lavel
          $('#saved-label-' + $termId + "-" + $lang)
              .transition('scale')
              .transition('scale');

          // if there is some warning mark "no default translation" then remove it
          if ($isDefault) {
            $('#warn-icon-' + $termId + '-def-lang').hide();
            $('#warn-label-def-lang-' + $termId).hide();
          }

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

        $('#table-terms').append(
          '<div class="item" id="row-' + data.term.ID + '">' +
            '<div class="content">' +
              '<a id="row-a-' + data.term.ID + '" class="header clickable-term" data-id="' + data.term.ID + '">' +
                  data.term.Code +
                  '<i id="warn-icon-' + data.term.ID + '-def-lang" class="warning sign icon withoutDefaultTranslation"></i>' +
              '</a>' +
          '</div>');

        $('#row-a-' + data.term.ID).click(fnClickTermCallback);

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
