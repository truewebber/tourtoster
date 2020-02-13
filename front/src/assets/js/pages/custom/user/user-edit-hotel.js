"use strict";

let UserEditHotel = function () {
    let formEl;

    let initSubmit = function () {
        let btn = formEl.find('input[type=submit]');

        btn.on('click', function (e) {
            e.preventDefault();

            // See: src\js\framework\base\app.js
            KTApp.progress(btn);
            //KTApp.block(formEl);

            $(formEl).ajaxSubmit({
                url: $(formEl).attr("action"),
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

    return {
        // public functions
        init: function () {
            formEl = $('form.new-user');

            initSubmit();
        }
    };
}();

jQuery(document).ready(function () {
    UserEditHotel.init();
});
