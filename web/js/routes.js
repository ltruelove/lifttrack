var app = Sammy('#main', function(){
    /***** HOME ROUTES *****/

    this.get('#/home', function(context){
        context.partial('/views/home.html', null, function(){
            var container = document.getElementById('homeForm')
            ko.cleanNode(container);
            ko.applyBindings(model,container);
        });
    });

    this.get('#/', function () {
        this.app.runRoute('get', '#/home')
    });

    this.post('#:id', function() {
        return false;
    });

    /***** END HOME ROUTES *****/

    /***** USER ROUTES *****/
    this.get('#/userHome', function(context){
        context.partial('/views/user/index.html', null, function(){
            var container = document.getElementById('userHome')
            ko.cleanNode(container);
            ko.applyBindings(model,container);
            model.getUserPrograms();
        });
    });

    this.get('#/user/new', function(context){
        context.partial('/views/user/new.html', null, function(){
            var container = document.getElementById('newUser')
            ko.cleanNode(container);
            ko.applyBindings(model,container);
            model.getNewUser();
        });
    });

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
    this.get('#/program/create', function(context){
        context.partial('/views/program/form.html',null, function(){
            var container = document.getElementById('programForm')
            ko.cleanNode(container);
            model.getProgram('0');
        });
    });
     
    this.get('#/program/:programId', function(context){
        var id = this.params['programId'];
        context.partial('/views/program/form.html',null, function(){
            var container = document.getElementById('programForm')
            ko.cleanNode(container);
            model.getProgram(id);
        });
    });
    /***** END PROGRAM ROUTES *****/
});
