$(document).ready(function() {

  // for the Docs, please refer to http://listjs.com/docs/
  var options = {
    valueNames: ['clickable-term'],
    multiSearch: false
  };
  var userList = new List('terms-list-panel', options);

  $('.ui.dropdown').dropdown();

});
