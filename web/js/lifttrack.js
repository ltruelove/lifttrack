var liftTrackViewModel = function(){
    var self = this;
    self.user = ko.mapping.fromJS({});
    self.token = "";
    self.loginUser = ko.mapping.fromJS({Username:"",Password:""});
    self.selectedProgram = ko.mapping.fromJS({});
    self.userPrograms = ko.mapping.fromJS([]);
    self.liftTypes = ko.mapping.fromJS([]);
    self.newUser = ko.mapping.fromJS({});

    self.login = function(){
        var loginUser = ko.mapping.toJS(self.loginUser);
        $.ajax({
            url: '/user/login',
            type: 'POST',
            data: JSON.stringify(loginUser),
            contentType: 'application/json',
            dataType: 'json',
            success: function(data, textStatus, request){
                self.token = request.getResponseHeader("Token");
                ko.mapping.fromJS(data, self.user);
                location.hash="/userHome";
            },
            error: function(request, textStatus, errorThrown){
                alert(errorThrown);
            }
        });
    };

    self.getUserPrograms = function(){
        $.ajax({
            url: '/user/programs',
            type: 'GET',
            dataType: 'json',
            beforeSend: function(request){
                request.setRequestHeader("Token", self.token)
            },
            success: function(data, textStatus, request){
                ko.mapping.fromJS(data, self.userPrograms);
            },
            error: function(request, textStatus, errorThrown){
                alert(errorThrown);
            }
        });
    };

    self.getProgram = function(id){
        $.ajax({
            url: '/program/' + id,
            type: 'GET',
            dataType: 'json',
            beforeSend: function(request){
                request.setRequestHeader("Token", self.token)
            },
            success: function(data, textStatus, request){
                ko.mapping.fromJS(data, self.selectedProgram);
            },
            error: function(request, textStatus, errorThrown){
                alert(errorThrown);
            }
        });
    };

    self.saveProgram = function(){
        var program = ko.mapping.toJS(self.selectedProgram);
        var type = (program.Id > 0)? 'PUT' : 'POST';
        var dataString = JSON.stringify(program);

        $.ajax({
            url: '/program',
            type: type,
            data: dataString,
            contentType: 'application/json',
            dataType: 'json',
            beforeSend: function(request){
                request.setRequestHeader("Token", self.token)
            },
            success: function(data, textStatus, request){
                ko.mapping.fromJS(data, self.selectedProgram);
                alert('Program Saved');
            },
            error: function(request, textStatus, errorThrown){
                alert(errorThrown);
            }
        });

    };

    self.addLiftToProgram = function(){
        $.ajax({
            url: '/lift/0',
            type: 'GET',
            dataType: 'json',
            beforeSend: function(request){
                request.setRequestHeader("Token", self.token)
            },
            success: function(data, textStatus, request){
                var program = ko.mapping.toJS(self.selectedProgram);
                //a blank lift does not contain the program id so we have to set it here
                data.ProgramId = program.Id;
                program.Lifts.push(data);
                ko.mapping.fromJS(program, self.selectedProgram);
            },
            error: function(request, textStatus, errorThrown){
                alert(errorThrown);
            }
        });
    };

    self.getLiftTypes = function(){
        $.get('/liftTypes',function(data){
            ko.mapping.fromJS(data, self.liftTypes);
        },'json');
    };

    self.getNewUser = function(){
        $.ajax({
            url: '/user/0',
            type: 'GET',
            dataType: 'json',
            success: function(data, textStatus, request){
                ko.mapping.fromJS(data, self.newUser);
            },
            error: function(request, textStatus, errorThrown){
                alert(errorThrown);
            }
        });
    };

    self.saveNewUser = function(){
        if($('#shippingAddress').val() != ""){
            //someone put their hand in the honey pot
            location.hash = '/home';
            return;
        }

        //@todo: get some better validation working in here
        if($('#passwordConfirm').val() != self.newUser.Password()){
            alert("Passwords don't match");
            return;
        }

        var user = ko.mapping.toJS(self.newUser);
        var dataString = JSON.stringify(user);

        $.ajax({
            url: '/user',
            type: 'POST',
            data: dataString,
            contentType: 'application/json',
            dataType: 'json',
            success: function(data, textStatus, request){
                ko.mapping.fromJS({}, self.newUser);
                alert('User Saved');
                location.hash = '/home';
            },
            error: function(request, textStatus, errorThrown){
                alert(errorThrown);
            }
        });

    };
};

var model = new liftTrackViewModel();
model.getLiftTypes();

$().ready(function(){
    app.run('#/');

});
