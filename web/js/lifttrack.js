var liftTrackViewModel = function(){
    var self = this;
    self.user = ko.mapping.fromJS({});
    self.selectedProgram = ko.mapping.fromJS({});

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
