function AppViewModel() {
    const self = this;

    self.email = ko.observable("");
    self.password = ko.observable("");
    self.token = ko.observable(localStorage.getItem("token") || "");
    self.error = ko.observable("");
    self.message = ko.observable("Welcome to the secured app!");

    self.isAuthenticated = ko.computed(() => !!self.token());

    self.login = function () {
        fetch("/api/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: self.email(),
                password: self.password()
            })
        })
        .then(res => {
            if (!res.ok) throw new Error("Login failed");
            return res.json();
        })
        .then(data => {
            self.token(data.token);
            localStorage.setItem("token", data.token);
            self.error("");
        })
        .catch(err => self.error(err.message));
    };

    self.logout = function () {
        self.token("");
        localStorage.removeItem("token");
        self.email("");
        self.password("");
    };
}

ko.applyBindings(new AppViewModel());