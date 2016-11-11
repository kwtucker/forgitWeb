var currentArticle = "tab_1",
  currentTab = 1,
  edit = 0;
$(function() {
  // Show hide form when not edit
  $('#newFormBody').hide()
  $('#formBody').hide()
  $('#editButton').on('click', function(){
    $('#formOverview').hide()
    $('#formBody').show()
    $('#formBody').css({
      "border": "2px solid #535E62",
      "border-top": 0,
      "padding": "0 20px 20px"
    })
  })

  $('#newForm').on('click', function(){
    $('#formOverview').hide()
    $('#formBody').hide()
    $('#newFormBody').show()
    $('#newFormBody').css({
      "border": "2px solid #535E62",
      "border-top": 0,
      "padding": "0 20px 20px"
    })
    $('#newFormBody input[name="settingGroupName"]').focus()
  });

  $('#cancelNewSetting').on('click', function(){
    $('#newFormBody').hide();
    $('#formOverview').show();
    window.scrollTo(0, 750);
  })

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

  // When a link on the subnav is selected it will add styling as activePage.
  var currentPage;
  switch (window.location.pathname) {
    case "/getting-started/":
    $("a[href='"+currentHash+"']").removeClass("activePage");
      $("a[href='/getting-started/']").addClass("activePage")
      currentPage = '/getting-started/';
      break;
    case "/dashboard/":
    $("a[href='"+currentHash+"']").removeClass("activePage");
      $("a[href='/dashboard/']").addClass("activePage")
      currentPage = '/dashboard/';
      break;
    default:
      $("a[href="+currentPage +"]").removeClass("activePage")
  }

  // When a link on a current page that does not go to another
  // page. It will have activePage styling for the current hash.
  var currentHash = "/";
  $("a[href='/#features']").on('click', function() {
    $("a[href='"+currentHash+"']").removeClass("activePage");
    $("a[href='/#features']").addClass("activePage");
    currentHash = '/#features';
  })
  $("a[href='/#price']").on('click', function() {
    $("a[href='"+currentHash+"']").removeClass("activePage");
    $("a[href='/#price']").addClass("activePage");
    currentHash = '/#price';
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

// Parse url then notify submit, remove or logout
(function (window, undefined) {
    "use strict";
    var urlParams;
    // onpopstate requires a function call
    (window.onpopstate = function () {
        var match,
            pl = /\+/g, // Regex for replacing addition symbol with a space
            // ([^&=]+) First group matches everything literally that is not &. "+" 1 match to unlimited.
            // ([^&]*)  ignores the & on everything. /g means global.
            search = /([^&=]+)=?([^&]*)/g,
            decode = function (s) {
              // Replace any + with space
              return decodeURIComponent(s.replace(pl, " "));
            },
            // ignore the ? and start with 2nd character in string.
            query = window.location.search.substring(1);
            // empty obj
            urlParams = {};
        // execute regex expression for each match and
        while (match = search.exec(query))
            urlParams[decode(match[1])] = decode(match[2]);
    })();

    // Style for Form Submit or log out with styling. Notify plugin
    $.notify.addStyle('SuccessfullySubmit', {
      html: "<div><span data-notify-text/></div>",
      classes: {
        base: {
          "white-space": "nowrap",
          "background-color": "#4DC4AD",
          "padding": "20px",
          "color": "#393939",
          "font": "500 1.3em 'Ubuntu', sans-serif",
          "min-width": "400px",
          "text-align": "center",
        },
      }
    });

    // Style for Form remove with styling. Notify plugin
    $.notify.addStyle('RemovedSetting', {
      html: "<div><span data-notify-text/></div>",
      classes: {
        base: {
          "white-space": "nowrap",
          "background-color": "#BE5D59",
          "padding": "20px",
          "color": "#E5E5E5",
          "font": "500 1.3em 'Ubuntu', sans-serif",
          "width": "400px",
          "text-align": "center",
        },
      }
    });

    // If the url param s is true the display setting submitted
    if (urlParams.s == "true") {
      $.notify('Setting Group Submitted', {
        style: 'SuccessfullySubmit'
      });
    }
    // If the url param r is true the display setting removed
    if (urlParams.r == "true") {
      $.notify('Setting Group Removed', {
        style: 'RemovedSetting'
      });
    }

    // Logout notify
    if (urlParams.lo == "true") {
      $.notify('Log Out Successful. See you next time.', {
        style: 'SuccessfullySubmit'
      });
    }

}) (window);
