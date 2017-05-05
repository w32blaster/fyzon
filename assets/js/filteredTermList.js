$(document).ready(function() {

  var options = {
    valueNames: ['clickable-term']
  };

  var userList = new List('terms-list-panel', options);

  $('.ui.dropdown').dropdown();

});
