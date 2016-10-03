$(function() {
    $("#tabs a").addClass('hov');
    $("#tabs a[data-tab='1']").removeClass('hov').addClass('selectedtab')
    $(".tab_1").css('display','flex');

    $('#tabs a').on('click', function() {
      if ($(this).hasClass('selectedtab') === false) {
        $('.selectedtab').addClass('hov');
        $('.selectedtab').removeClass('selectedtab');
        $(this).addClass('selectedtab');
      }
      $('#features article').hide()
      $(this).addClass("selectedtab")
      $(this).removeClass('hov');
      var tab = $(this).data('tab')
      $('.tab_'+tab).css('display','flex')
    })
});
