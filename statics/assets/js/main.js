$(function() {
    console.log("Pitel");
    try {
        $(".js-select2").select2({
            closeOnSelect: false,
            placeholder: "Chọn các thể loại",
            allowHtml: true,
            allowClear: true,
            tags: true // ÑÐ¾Ð·Ð´Ð°ÐµÑ‚ Ð½Ð¾Ð²Ñ‹Ðµ Ð¾Ð¿Ñ†Ð¸Ð¸ Ð½Ð° Ð»ÐµÑ‚Ñƒ
        });
        $(".js-select2").on("change", function() {
            let val = $(this).val()
            $("input[name='category']").val(val)
        })

        $('.icons_select2').select2({
            width: "100%",
            templateSelection: iformat,
            templateResult: iformat,
            allowHtml: true,
            placeholder: "Click to select an option",
            dropdownParent: $('.select-icon'), //Ð¾Ð±Ð°Ð²Ð¸Ð»Ð¸ ÐºÐ»Ð°ÑÑ
            allowClear: true,
            multiple: false
        });
    } catch (error) {

    }

    function extractHostname(url) {
        var hostname;
        if (url.indexOf("//") > -1) {
            hostname = url.split('/')[2];
        } else {
            hostname = url.split('/')[0];
        }
        hostname = hostname.split(':')[0];
        hostname = hostname.split('?')[0];
        return hostname;
    }
    $('#url').on('blur', function(e) {
        let val = $(this).val()
        let hn = extractHostname(val);
        $("select[name='source']")[0].value = hn;
        console.log(val)
    });

    function iformat(icon, badge, ) {
        var originalOption = icon.element;
        var originalOptionBadge = $(originalOption).data('badge');

        return $('<span><i class="fa ' + $(originalOption).data('icon') + '"></i> ' + icon.text + '<span class="badge">' + originalOptionBadge + '</span></span>');
    }

    // $("#category_select").on("change", function(e) {
    //     e.preventDefault();
    //     let val = $(this).val()
    //     $("input[name='category']").val(val)
    // });
    $("a[name='delete']").on("click", function(e) {
        e.preventDefault();
        let val = $(this).attr('tag')
        swal({
                title: "Are you sure?",
                text: "Once deleted, you will not be able to recover this url",
                icon: "warning",
                buttons: true,
                dangerMode: true,
            })
            .then((willDelete) => {
                if (willDelete) {
                    window.location.href = val;
                } else {
                    swal("Canceled!");
                }
            });
    });
    $("#get_preview").on("click", function(e) {
        e.preventDefault();
        var url = $("input[name='url']")[0].value;
        let source = $("select[name='source']")[0].value;
        if (url == "") {
            swal("Lỗi", "URL rỗng")
            return;
        }
        swal("Đang lấy thông tin", "Vui lòng đợi")
        var settings = {
            "url": `/preview`,
            "method": "POST",
            "timer": 0,
            "headers": {
                "Content-Type": "application/json"
            },
            "data": JSON.stringify({
                "url": `${url}`,
                "source": `${source}`,
            }),
        };
        $.ajax(settings).fail(function(response) {
            swal("Lỗi", "Không thể lấy thông tin");
        }).done(function(response) {
            console.log(response)
            if (response.code != 200) {
                swal("Lỗi", "Không thể lấy thông tin");
                return;
            }

            if (response.info == undefined) {
                swal("Thông báo", "Không có thông tin");
                return;
            }
            $("#information")[0].innerHTML =
                `
            <div class="col-md-6">
            <h4>Thông tin mô tả</h4>
            ${response.info.replaceAll("h3", "h5").replaceAll('<a', '<a disabled ')}
            </div>
            
            <div class="col-md-6">
            <img src='${response.image}'>
            </div>
            
            `;
            $("#information a").addClass("disabled")
        });

    });
    // $("select[name='category_select'] > option")[0].selected = true
    // $("select[name='category_select']").trigger('change');
});