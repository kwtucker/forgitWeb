var currentArticle = "tab_1"
var currentTab = 1
var edit = 0
$(function() {
  // Show hide form when not edit
  $('#newFormBody').hide()
  $('#formBody').hide()
  $('#editButton').on('click', function(){
    $('#formOverview').hide()
    $('#clearForm').show()
    $('#formBody').show()
    $('#formBody').css({
      "border": "2px solid #535E62",
      "border-top": 0,
      "padding": "0 20px 20px"
    })
  })

  $('#newForm').on('click', function(){
    $('#formOverview').hide()
    $('#newFormBody').show()
    $('#newFormBody').css({
      "border": "2px solid #535E62",
      "border-top": 0,
      "padding": "0 20px 20px"
    })
  });


  if ($(window).width() < 1000) {
    $('#features article').css('display','flex')
    $('#features article:even').addClass('even')
    $('#features').css('background', '#4DC4AD')
  } else {
    $('#features article').hide()
    $('#features article').removeClass('even')
    $('#tabs a').removeClass('selectedtab')
    $('#tabs a[data-tab='+currentTab+']').addClass('selectedtab')
    $('.tab_'+currentTab).css('display','flex')
    $('#features').css('background', '#5E696D')
  }

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
    currentTab = tab
  })
});

$(window).resize(function() {
  if ($(window).width() < 1000) {
    $('#features article').css('display','flex')
    $('#features article:even').addClass('even')
    $('#features').css('background', '#4DC4AD')
  } else {
    $('#features article').hide()
    $('#features article').removeClass('even')
    $('#tabs a').removeClass('selectedtab')
    $('#tabs a[data-tab='+currentTab+']').addClass('selectedtab')
    $('.tab_'+currentTab).css('display','flex')
    $('#features').css('background', '#5E696D')
  }
});
