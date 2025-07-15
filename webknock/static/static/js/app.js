function ViewModel() {
    var self = this;
    self.currentView = ko.observable('login');
    self.loginUsername = ko.observable('');
    self.loginPassword = ko.observable('');
    self.registerUsername = ko.observable('');
    self.registerPassword = ko.observable('');
    self.books = ko.observableArray([]);
    self.members = ko.observableArray([]);
    self.newBookTitle = ko.observable('');
    self.newBookAuthor = ko.observable('');
    self.newMemberName = ko.observable('');
    self.editingBook = ko.observable(null);
    self.editingMember = ko.observable(null);
    self.filterStatus = ko.observable('');
    self.selectedBook = ko.observable();
    self.selectedMember = ko.observable();
    self.token = ko.observable(localStorage.getItem('token') || '');

    self.switchToRegister = function() {
        self.currentView('register');
    };

    self.switchToLogin = function() {
        self.currentView('login');
    };

    self.switchToBooks = function() {
        self.currentView('books');
        self.loadBooks();
        self.loadMembers();
    };

    self.switchToMembers = function() {
        self.currentView('members');
        self.loadBooks();
        self.loadMembers();
    };

    self.login = function() {
        axios.post('/api/login', {
            username: self.loginUsername(),
            password: self.loginPassword()
        }).then(function(response) {
            self.token(response.data.token);
            localStorage.setItem('token', response.data.token);
            self.currentView('books');
            self.loadBooks();
            self.loadMembers();
        }).catch(function(error) {
            alert('Login failed: ' + error.response.data);
        });
    };

    self.register = function() {
        axios.post('/api/register', {
            username: self.registerUsername(),
            password: self.registerPassword()
        }).then(function() {
            alert('Registration successful');
            self.switchToLogin();
        }).catch(function(error) {
            alert('Registration failed: ' + error.response.data);
        });
    };

    self.logout = function() {
        self.token('');
        localStorage.removeItem('token');
        self.currentView('login');
        self.books([]);
        self.members([]);
    };

    self.loadBooks = function() {
        axios.get('/api/books', {
            headers: { Authorization: 'Bearer ' + self.token() },
            params: { status: self.filterStatus() }
        }).then(function(response) {
            self.books(response.data);
        }).catch(function(error) {
            alert('Error loading books: ' + error.response.data);
        });
    };

    self.loadMembers = function() {
        axios.get('/api/members', {
            headers: { Authorization: 'Bearer ' + self.token() }
        }).then(function(response) {
            self.members(response.data);
        }).catch(function(error) {
            alert('Error loading members: ' + error.response.data);
        });
    };

    self.saveBook = function() {
        var book = {
            title: self.newBookTitle(),
            author: self.newBookAuthor(),
            status: 'available'
        };
        if (self.editingBook()) {
            axios.put('/api/books/' + self.editingBook().id, book, {
                headers: { Authorization: 'Bearer ' + self.token() }
            }).then(function() {
                self.editingBook(null);
                self.newBookTitle('');
                self.newBookAuthor('');
                self.loadBooks();
            }).catch(function(error) {
                alert('Error updating book: ' + error.response.data);
            });
        } else {
            axios.post('/api/books', book, {
                headers: { Authorization: 'Bearer ' + self.token() }
            }).then(function() {
                self.newBookTitle('');
                self.newBookAuthor('');
                self.loadBooks();
            }).catch(function(error) {
                alert('Error adding book: ' + error.response.data);
            });
        }
    };

    self.editBook = function(book) {
        self.editingBook(book);
        self.newBookTitle(book.title);
        self.newBookAuthor(book.author);
    };

    self.deleteBook = function(book) {
        axios.delete('/api/books/' + book.id, {
            headers: { Authorization: 'Bearer ' + self.token() }
        }).then(function() {
            self.loadBooks();
        }).catch(function(error) {
            alert('Error deleting book: ' + error.response.data);
        });
    };

    self.saveMember = function() {
        var member = { name: self.newMemberName() };
        if (self.editingMember()) {
            axios.put('/api/members/' + self.editingMember().id, member, {
                headers: { Authorization: 'Bearer ' + self.token() }
            }).then(function() {
                self.editingMember(null);
                self.newMemberName('');
                self.loadMembers();
            }).catch(function(error) {
                alert('Error updating member: ' + error.response.data);
            });
        } else {
            axios.post('/api/members', member, {
                headers: { Authorization: 'Bearer ' + self.token() }
            }).then(function() {
                self.newMemberName('');
                self.loadMembers();
            }).catch(function(error) {
                alert('Error adding member: ' + error.response.data);
            });
        }
    };

    self.editMember = function(member) {
        self.editingMember(member);
        self.newMemberName(member.name);
    };

    self.deleteMember = function(member) {
        axios.delete('/api/members/' + member.id, {
            headers: { Authorization: 'Bearer ' + self.token() }
        }).then(function() {
            self.loadMembers();
        }).catch(function(error) {
            alert('Error deleting member: ' + error.response.data);
        });
    };

    self.borrowBook = function() {
        axios.post('/api/borrow', {
            book_id: self.selectedBook(),
            member_id: self.selectedMember()
        }, {
            headers: { Authorization: 'Bearer ' + self.token() }
        }).then(function() {
            self.loadBooks();
            self.loadMembers();
        }).catch(function(error) {
            alert('Error borrowing book: ' + error.response.data);
        });
    };

    self.returnBook = function() {
        axios.post('/api/return', {
            book_id: self.selectedBook(),
            member_id: self.selectedMember()
        }, {
            headers: { Authorization: 'Bearer ' + self.token() }
        }).then(function() {
            self.loadBooks();
            self.loadMembers();
        }).catch(function(error) {
            alert('Error returning book: ' + error.response.data);
        });
    };

    if (self.token()) {
        self.switchToBooks();
    }
}

ko.applyBindings(new ViewModel());