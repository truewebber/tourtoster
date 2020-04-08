/*
	Add image by link plugin for Imperavi Redactor v10.x.x
	Updated: January 21, 2015

	Copyright (c) 2015, Alexey Kinev
	License: The MIT License (MIT) http://opensource.org/licenses/MIT

	Usage:
        	$('#redactor').redactor({
        		imageUpload: false,
        		s3: false,
        		plugins: ['imagelink']
        	});
*/
if (!RedactorPlugins) var RedactorPlugins = {};

(function($)
{
    RedactorPlugins.imagelink = function()
    {
        return {
            getTemplate: function()
            {
                return String()
                    + '<section id="redactor-modal-imagelink-insert">'
                    + '<label>' + this.lang.get('image_web_link') + '</label>'
                    + '<input type="text" id="redactor-insert-imagelink">'
                    + '</section>';
            },
            init: function()
            {
                var button = this.button.addAfter('underline', 'image', this.lang.get('image'));
                this.button.addCallback(button, this.imagelink.show);
            },
            show: function()
            {
                this.modal.addTemplate('imagelink', this.imagelink.getTemplate());

                this.modal.load('imagelink', this.lang.get('image'), 700);
                this.modal.createCancelButton();

                var button = this.modal.createActionButton(this.lang.get('insert'));
                button.on('click', this.imagelink.insert);

                this.selection.save();
                this.modal.show();

                $('#redactor-insert-imagelink').focus();
            },
            insert: function()
            {
                var data = $('#redactor-insert-imagelink').val();
                this.image.insert({filelink: data}, false);
            }
        };
    };
})(jQuery);
