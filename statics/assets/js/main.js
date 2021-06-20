$(function() {
    console.log("Pitel");
    $("#category_select").on("change", function(e) {
        e.preventDefault();
        let val = $(this).val()
        $("input[name='category']").val(val)
    });
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
            ${response.info.replaceAll("h3","h5").replaceAll('<a','<a disabled ')}
            </div>
            
            <div class="col-md-6">
            <img src='${response.image}'>
            </div>
            
            `;
            $("#information a").addClass("disabled")
        });

    });
    
    $("select[name='category_select'] > option")[0].selected = true
    $("select[name='category_select']").trigger('change');
});