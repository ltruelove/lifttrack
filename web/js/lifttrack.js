var liftTrackViewModel = function(){
    var self = this;
    self.user = ko.mapping.fromJS({});
    self.token = "";
    self.loginUser = ko.mapping.fromJS({Username:"",Password:""});
    self.selectedProgram = ko.mapping.fromJS({});
    self.userPrograms = ko.mapping.fromJS([]);

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

$().ready(function(){
    app.run('#/');

});
