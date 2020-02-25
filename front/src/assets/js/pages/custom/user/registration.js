"use strict";

let Registration = function () {
    let registrationFormEl;

    let showErrorMsg = function (form, type, msg) {
        var alert = $('<div class="alert alert-' + type + ' alert-dismissible" role="alert">\
			<div class="alert-text">' + msg + '</div>\
			<div class="alert-close">\
                <i class="flaticon2-cross kt-icon-sm" data-dismiss="alert"></i>\
            </div>\
		</div>');

        form.find('.alert').remove();
        alert.prependTo(form);
        //alert.animateClass('fadeIn animated');
        KTUtil.animateClass(alert[0], 'fadeIn animated');
        alert.find('span').html(msg);
    };

    let initSubmit = function () {
        let btn = registrationFormEl.find('input[type=submit]');

        btn.on('click', function (e) {
            e.preventDefault();

            registrationFormEl.validate({
                rules: {
                    first_name: {
                        required: true
                    },
                    last_name: {
                        required: true
                    },
                    email: {
                        required: true,
                        email: true
                    },
                    phone: {
                        required: true
                    },
                    password: {
                        required: true
                    },
                    password_repeat: {
                        required: true
                    },
                }
            });

            if (!registrationFormEl.valid()) {
                return;
            }

            btn.addClass('kt-spinner kt-spinner--right kt-spinner--sm kt-spinner--light').attr('disabled', true);

            registrationFormEl.ajaxSubmit({
                url: $(registrationFormEl).attr("action"),
                type: "post",
                dataType: "json",
                success: function () {
                    btn.removeClass('kt-spinner kt-spinner--right kt-spinner--sm kt-spinner--light').attr('disabled', false); // remove
                    registrationFormEl.clearForm(); // clear form
                    registrationFormEl.validate().resetForm(); // reset validation states

                    showErrorMsg(registrationFormEl, 'success', 'Cool! We will contact you soon.');
                },
                error: function (r) {
                    btn.removeClass('kt-spinner kt-spinner--right kt-spinner--sm kt-spinner--light').attr('disabled', false); // remove

                    let errorText = "Unknown server error";
                    if (r.responseJSON) {
                        errorText = r.responseJSON.error;
                    }

                    showErrorMsg(registrationFormEl, 'danger', errorText);
                }
            })
        });
    };

    return {
        // public functions
        init: function () {
            registrationFormEl = $('form.registration-from');

            initSubmit();
        },
    };
}();

jQuery(document).ready(function () {
    Registration.init();
});
