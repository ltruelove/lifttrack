var liftTrackViewModel = function(){
    var self = this;
    self.user = ko.mapping.fromJS({});
    self.token = "";
    self.loginUser = ko.mapping.fromJS({Username:"",Password:""});
    self.selectedProgram = ko.mapping.fromJS({});
    self.userPrograms = ko.mapping.fromJS([]);
    self.liftTypes = ko.mapping.fromJS([]);

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
    };

    self.getLiftTypes = function(){
        $.get('/liftTypes',function(data){
            ko.mapping.fromJS(data, self.liftTypes);
        },'json');
    }
    /*
    self.getInventoryItems = function(container){
        $.get('/inventory', function(data){
            self.inventoryItems(data);
            ko.applyBindings(self,container);
        },'json');
    };

    self.saveAuction = function(){
        var auction = ko.mapping.toJS(self.selectedAuction);
        var type = (auction.Id > 0)? 'PUT':'POST';

        $.ajax({
            url: '/auction',
            contentType: 'application/json',
            dataType: 'json',
            type: type,
            data: JSON.stringify(auction),
            success: function(data) {
                ko.mapping.fromJS(data, self.selectedAuction);
                alert('Auction saved');
                location.hash = "/inventory/auctions";
            },
            error: function(request, status, errorString) {
                alert(errorString);
            },
            cache: false
          }
       );
    };
    */

};

var model = new liftTrackViewModel();
model.getLiftTypes();

$().ready(function(){
    app.run('#/');

});
