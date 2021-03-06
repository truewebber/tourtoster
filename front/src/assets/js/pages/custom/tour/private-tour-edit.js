"use strict";

import {RRule, RRuleSet} from 'rrule';
import {Calendar} from '@fullcalendar/core';
import dayGridPlugin from '@fullcalendar/daygrid';
import rrulePlugin from '@fullcalendar/rrule';

let PrivateTourEdit = function () {
    let tourFormEl;
    let rrEl;
    let rrFormEl;
    let rrFormToogle;
    let rruleSet;
    let calendarEl;
    let textRRuleSetEl;
    let calendar;
    let imageUploadEl;
    let imageTarget;

    let parseIntArray = function (str) {
        let intArr = [];
        str = str.trim(" ")
        if (str !== "") {
            str.split(",").forEach(elem => intArr.push(parseInt(elem)));
        }

        return intArr;
    }

    let clearRRuleForm = function () {
        rrFormEl.find('input[type!=checkbox][type!=radio]').val('');
        rrFormEl.find('input[type=checkbox]').prop('checked', false);
        rrFormEl.find('input[type=radio]').prop('checked', false);
    }

    let closeRRuleForm = function () {
        rrFormToogle.text('Add new rule').attr({'data-action': 'closed'});
        rrFormEl.hide();
    }

    let openRRuleForm = function () {
        rrFormToogle.text('Cancel').attr({'data-action': 'open'});
        rrFormEl.show();
    }

    let drawRules = function () {
        let rrules = rruleSet.rrules();
        let exrules = rruleSet.exrules();

        if (rrules.length === 0 && exrules.length === 0) {
            rrEl.html('No rules here yet.');
            renderCalendar();
            textRRuleSet();

            return;
        }

        rrEl.empty();

        if (rrules.length !== 0) {
            rrEl.append('<h3>rrules</h3>');
            for (let i = 0; i < rrules.length; i++) {
                let close = $('<span>').html('x');
                let r = rrules[i].toString().replace('RRULE:', '');
                let elem = $('<p>').attr({'data-type': '1', 'data-index': i}).html(r).append(close);

                rrEl.append(elem);
            }
        }

        if (exrules.length !== 0) {
            rrEl.append('<h3>exrules</h3>');
            for (let i = 0; i < exrules.length; i++) {
                let close = $('<span>').html('x');
                let r = exrules[i].toString().replace('RRULE:', '');
                let elem = $('<p>').attr({'data-type': '2', 'data-index': i}).html(r).append(close);

                rrEl.append(elem);
            }
        }

        initDeleteRules();
        renderCalendar();
        textRRuleSet();
    }

    let initDeleteRules = function () {
        let buttons = rrEl.find('p > span');

        $(buttons).on('click', function () {
            let elem = $(this).parent();
            let index = parseInt($(elem).attr('data-index'));

            switch ($(elem).attr('data-type')) {
                case '1':
                    if (index < 0 || index > rruleSet._rrule.length - 1) {
                        return;
                    }

                    delete rruleSet._rrule[index];
                    rruleSet._rrule.length--;
                    break;
                case '2':
                    if (index < 0 || index > rruleSet._exrule.length - 1) {
                        return;
                    }

                    delete rruleSet._exrule[index];
                    rruleSet._exrule.length--;
                    break;
            }

            drawRules();
        });
    }

    let initToggleRRuleForm = function () {
        rrFormToogle.on('click', function (e) {
            e.preventDefault();

            let current = $(this).attr('data-action');
            if (current === undefined) {
                current = "closed"
            }

            if (current === "closed") {
                openRRuleForm();
            } else if (current === "open") {
                clearRRuleForm();
                closeRRuleForm();
            }
        });
    }

    let initAddRule = function () {
        let btn = rrFormEl.find('button');

        btn.on('click', function (e) {
            e.preventDefault();

            let type = $('input[name=rrule_type]:checked').val();
            if (type === undefined) {
                alert('set type!');
                return
            }

            let freq = $('input[name=rrule_freq]:checked').val();
            if (freq === undefined) {
                alert('set frequency!');
                return
            }

            let byweekday = $('input[name=rrule_byweekday]:checked').map(function () {
                return parseInt($(this).val())
            }).get();
            let bymonth = $('input[name=rrule_bymonth]:checked').map(function () {
                return parseInt($(this).val())
            }).get();

            let until = $('input[name=rrule_until]').val();
            let count = $('input[name=rrule_count]').val();
            let interval = $('input[name=rrule_interval]').val();

            let bymonthday = parseIntArray($('input[name=rrule_bymonthday]').val());
            let byyearday = parseIntArray($('input[name=rrule_byyearday]').val());
            let byweekno = parseIntArray($('input[name=rrule_byweekno]').val());

            let options = {
                freq: freq,
                byweekday: byweekday,
                bymonth: bymonth,
                bymonthday: bymonthday,
                byyearday: byyearday,
                byweekno: byweekno,
            };

            if (until !== "") {
                options.until = until;
            }
            if (count !== "") {
                options.count = count;
            }
            if (interval !== "") {
                options.interval = interval;
            }

            let rrule = new RRule(options);

            switch (type) {
                case '1':
                    rruleSet.rrule(rrule);
                    break;
                case '2':
                    rruleSet.exrule(rrule);
                    break;
                default:
                    alert('unknown type');
                    break;
            }

            clearRRuleForm();
            drawRules();
        });
    }

    let textRRuleSet = function () {
        $(textRRuleSetEl).val(rruleSet.toString());
    }

    let renderCalendar = function () {
        let events = []

        if (rruleSet.toString() !== '') {
            let currentDate = new Date();
            let tomorrowDate = new Date();
            tomorrowDate.setDate(tomorrowDate.getDate());

            let tmpRS = rruleSet.clone();
            tmpRS.exdate(currentDate);
            tmpRS.exdate(tomorrowDate);

            events.push({
                id: 'onlyEvent',
                allDay: true,
                rrule: tmpRS.toString(),
                display: 'background'
            })
        }

        calendar = new Calendar(calendarEl, {
            plugins: [dayGridPlugin, rrulePlugin],
            initialView: 'dayGridMonth',
            height: 'auto',
            events: events
        });

        calendar.render();
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
    }

    let initLoadImage = function () {
        $(imageUploadEl).on('change', function () {
            let input = $(this);
            let files = input.prop('files');

            var reader = new FileReader();
            reader.onload = function (e) {
                $(imageTarget).css({'background-image': 'url(' + e.target.result + ')'}).removeClass('no-image');

            };

            reader.readAsDataURL(files[0]);
        });
    }

    return {
        // public functions
        init: function () {
            tourFormEl = $('form.new-tour');
            rrEl = $('div.recurrence-rules');
            rrFormEl = $('div.recurrence-rule-set');
            rrFormToogle = $('.recurrence-rule-set-toogle');
            calendarEl = $('div.recurrence-rule-calendar')[0];
            textRRuleSetEl = $('input[name=rrule_set]')[0];
            rruleSet = new RRuleSet();
            imageUploadEl = $('#file_image')[0];
            imageTarget = $('.tour-image-block > .preview');

            initSubmitTour();
            initToggleRRuleForm();
            initAddRule();
            initLoadImage();
            renderCalendar();
            textRRuleSet();
        },
    };
}();

jQuery(document).ready(function () {
    PrivateTourEdit.init();
});
