"use strict";

let UserEditHotel = function () {
    let userFormEl;
    let hotelFormEl;

    let hotelPopupEl;
    let hotelListEl;

    let initSubmitUser = function () {
        let btn = userFormEl.find('input[type=submit]');

        btn.on('click', function (e) {
            e.preventDefault();

            // See: src\js\framework\base\app.js
            KTApp.progress(btn);
            //KTApp.block(formEl);

            $(userFormEl).ajaxSubmit({
                url: $(userFormEl).attr("action"),
                type: "post",
                dataType: "json",
                success: function () {
                    KTApp.unprogress(btn);
                    // KTApp.unblock(formEl);

                    swal.fire({
                        "title": "",
                        "text": "The application has been successfully submitted!",
                        "type": "success",
                        "confirmButtonClass": "btn btn-secondary"
                    });
                },
                error: function (r) {
                    KTApp.unprogress(btn);

                    let errorText = "Unknown server error";
                    if (r.responseJSON) {
                        errorText = r.responseJSON.error;
                    }

                    swal.fire({
                        "title": "error",
                        "text": errorText,
                        "type": "error",
                        "confirmButtonClass": "btn btn-secondary"
                    });
                }
            })
        });
    };

    let initOpenPopupHotel = function () {
        let btn = $('#popup-hotel-open');

        $(btn).on('click', function (e) {
            e.preventDefault();

            loadHotelList();

            $(hotelPopupEl).show();
            $('body').addClass('body-popup-overflow');
        });
    };

    let initClosePopupHotel = function () {
        let btn = $('#popup-hotel-close');

        $(btn).on('click', function (e) {
            e.preventDefault();

            $(hotelPopupEl).hide();
            $('body').removeClass('body-popup-overflow');
        });
    };

    let initSubmitHotel = function () {
        let btn = hotelFormEl.find('input[type=submit]');

        btn.on('click', function (e) {
            e.preventDefault();

            $(hotelFormEl).ajaxSubmit({
                url: $(hotelFormEl).attr("action"),
                type: "post",
                dataType: "json",
                success: function () {
                    loadHotelList();
                },
                error: function (r) {
                    alert(r.responseJSON.error);
                }
            })
        });
    };

    let initNewHotel = function () {
        let btn = $('.popup .new-elem a');

        btn.on('click', function (e) {
            e.preventDefault();

            hotelFormEl.find('input[name=id]').val("");
            hotelFormEl.find('input[name=name]').val("");
        });
    };

    let initEditHotel = function () {
        let btn = hotelListEl.find('.elem .edit');

        btn.on('click', function (e) {
            e.preventDefault();

            hotelFormEl.find('input[name=id]').val($(this).attr("data-id"));
            hotelFormEl.find('input[name=name]').val($(this).attr("data-name"));
        });
    };

    let initRemoveHotel = function () {
        let btn = hotelListEl.find('.elem .remove');

        btn.on('click', function (e) {
            e.preventDefault();

            if (!confirm("Are you sure?")) {
                return
            }

            hotelFormEl.find('input[name=id]').val("");
            hotelFormEl.find('input[name=name]').val("");

            $.ajax({
                url: "/console/api/hotel",
                type: "delete",
                data: "id=" + $(this).attr("data-id"),
                dataType: "json",
                success: function (r) {
                    loadHotelList();
                },
                error: function (r) {
                    alert(r.responseJSON.error);

                    console.log("error delete hotel", r);
                }
            })
        });
    };

    let loadHotelList = function () {
        $(hotelListEl).html("Loading...");

        $.ajax({
            url: "/console/api/hotel",
            type: "get",
            dataType: "json",
            success: function (r) {
                $(hotelListEl).html("");

                r.forEach(function (e) {
                    $(hotelListEl).append(
                        '<div class="elem">' +
                        '<div class="hotel_name">' + e.name + '</div>' +
                        '<a class="edit" data-id="' + e.id + '" data-name="' + e.name + '" href="javascript:void(0);">✎</a>' +
                        '<a class="remove" data-id="' + e.id + '" data-name="' + e.name + '" href="javascript:void(0);">✕</a>' +
                        '</div>'
                    );
                });

                initEditHotel();
                initRemoveHotel();
            },
            error: function (r) {
                $(hotelListEl).html("Error get hotel list.");

                console.log("error get hotel list", r);
            }
        })
    };

    return {
        // public functions
        init: function () {
            userFormEl = $('form.new-user');
            hotelFormEl = $('form.hotel');
            hotelPopupEl = $('#popup-hotel');
            hotelListEl = $('#popup-hotel .list');

            initSubmitUser();
            initOpenPopupHotel();
            initClosePopupHotel();

            initSubmitHotel();
            initNewHotel();
        },
    };
}();

jQuery(document).ready(function () {
    UserEditHotel.init();
});
