"use strict";

let PrivateTourEdit = function () {
    let tourFormEl;
    let rrFormEl;

    let initToggleRRuleForm = function () {
        let btn = rrFormEl.find('input[type=submit]');


    }

    let initSubmitTour = function () {
        let btn = tourFormEl.find('input[type=submit]');

        btn.on('click', function (e) {
            e.preventDefault();

            tourFormEl.validate({
                rules: {
                    title: {
                        required: true
                    },
                    description: {
                        required: true
                    },
                },
                errorPlacement: function (error, element) {
                    error.css({"color": "#f00"});

                    if (element.is(":radio")) {
                        element = element.parent().parent();
                        element.append(error);

                        return;
                    }

                    error.insertAfter(element);
                },
            });

            if (!tourFormEl.valid()) {
                return;
            }

            KTApp.progress(btn);

            $(tourFormEl).ajaxSubmit({
                url: "/console/api/tour",
                type: "post",
                dataType: "json",
                success: function (r) {
                    KTApp.unprogress(btn);

                    swal.fire({
                        "title": "",
                        "text": "Private tour has been successfully submitted!",
                        "type": "success",
                        "confirmButtonClass": "btn btn-secondary"
                    }, function () {
                        window.location = "/console/private_tours/edit?id=" + r.responseJSON.id;
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
            tourFormEl = $('form.new-tour');
            rrFormEl = $('div.recurrence-rule-set');

            initSubmitTour();
        },
    };
}();

jQuery(document).ready(function () {
    PrivateTourEdit.init();
});
