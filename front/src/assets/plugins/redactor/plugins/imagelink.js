(function ($R) {
    $R.add('plugin', 'imagelink', {
        translations: {
            en: {
                "imagelink": "ImageLink",
                "image-link": "Insert image from link"
            }
        },
        modals: {
            'imagelink':
                '<form action=""> \
                    <div class="form-item"> \
                        <label for="modal-imagelink-input">## image-link ##</label> \
                        <input id="modal-imagelink-input" name="imagelink" /> \
                    </div> \
                </form>'
        },
        init: function (app) {
            this.app = app;
            this.lang = app.lang;
            this.opts = app.opts;
            this.toolbar = app.toolbar;
            this.component = app.component;
            this.insertion = app.insertion;
            this.inspector = app.inspector;
        },
        onmodal: {
            imagelink: {
                opened: function ($modal, $form) {
                    let IMGLNK = this

                    $form.getField('imagelink').focus();
                    $form.on('submit', function (e) {
                        e.preventDefault()

                        var data = $form.getData();
                        IMGLNK._insert(data);
                    })
                },
                insert: function ($modal, $form) {
                    var data = $form.getData();
                    this._insert(data);
                }
            }
        },
        oncontextbar: function (e, contextbar) {
            var data = this.inspector.parse(e.target)
            if (data.isComponentType('imagelink')) {
                var node = data.getComponent();
                var buttons = {
                    "remove": {
                        title: this.lang.get('delete'),
                        api: 'plugin.imagelink.remove',
                        args: node
                    }
                };

                contextbar.set(e, node, buttons, 'bottom');
            }

        },

        // public
        start: function () {
            var obj = {
                title: this.lang.get('imagelink'),
                api: 'plugin.imagelink.open'
            };

            var $button = this.toolbar.addButtonAfter('image', 'imagelink', obj);
            $button.setIcon('<i class="re-icon-image"></i>');
        },
        open: function () {
            var options = {
                title: this.lang.get('imagelink'),
                width: '600px',
                name: 'imagelink',
                handle: 'insert',
                commands: {
                    insert: {title: this.lang.get('insert')},
                    cancel: {title: this.lang.get('cancel')}
                }
            };

            this.app.api('module.modal.build', options);
        },
        remove: function (node) {
            this.component.remove(node);
        },

        // private
        _insert: function (data) {
            this.app.api('module.modal.close');

            if (data.imagelink.trim() === '') {
                return;
            }

            let d = new Date();
            let imageObj = {
                myimagekey: {
                    id: 'image' + d.toJSON(),
                    url: data.imagelink.trim()
                }
            };

            this.app.api('module.image.insert', imageObj);
        }
    });
})(Redactor);
