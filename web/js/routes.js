var app = Sammy('#main', function(){
    /***** HOME ROUTES *****/

    this.get('#/home', function(context){
        context.partial('/views/home.html');
    });

    this.get('#/', function () {
        this.app.runRoute('get', '#/home')
    });

    /***** END HOME ROUTES *****/

    /***** USER ROUTES *****/

    this.get('#/users', function(context){
        context.partial('/views/users.html', null, function(){
            var container = document.getElementById('userList')
            ko.cleanNode(container);
            ko.applyBindings(model,container);
            model.getUserList(container);
        });
    });

    this.get('#/user', function(context){
        context.partial('/views/users/info.html',null, function(){
            var container = document.getElementById('userInfo')
            ko.cleanNode(container);
            ko.applyBindings(model,container);
        });
    });

    this.get('#/user/:userId', function(context){
        var id = this.params['userId'];
        context.partial('/views/users/info.html',null, function(){
            var container = document.getElementById('userInfo')
            ko.cleanNode(container);
            model.fetchUser(container, id);
        });
    });

    /***** END USER ROUTES *****/

    /***** PROGRAM ROUTES *****/

    /***** END PROGRAM ROUTES *****/
});
