
const app = new (function () {
  const self = this;

  self.currentHtml = ko.observable('');
  self.currentViewModel = null;

  self.loadPage = function (pageName) {
    // Load HTML
    fetch(`/static/views/${pageName}.html`)
      .then(res => res.text())
      .then(html => {
        self.currentHtml(html);
        // Delay binding until HTML is inserted into DOM
        setTimeout(() => {
          self.loadViewModel(pageName);
          history.pushState({ page: pageName }, '', '/' + pageName);
        }, 0);
      });
  };

  self.loadViewModel = function (pageName) {
    // Unbind previous viewmodel if needed (manual teardown if necessary)
    if (self.currentViewModel && self.currentViewModel.dispose) {
      self.currentViewModel.dispose();
    }

    // Optional: Load JSON data
    fetch(`/static/data/${pageName}.json`)
      .then(res => res.json())
      .then(data => {
        // Create a new viewmodel based on data
        const viewModel = self.createViewModel(pageName, data);
        self.currentViewModel = viewModel;
        ko.applyBindings(viewModel, document.getElementById('view-container'));
      });
  };

  self.createViewModel = function (pageName, data) {
    switch (pageName) {
      case 'home':
        return {
          message: ko.observable(data.message)
        };
      case 'about':
        return {
          details: ko.observableArray(data.details)
        };
      default:
        return {};
    }
  };

  // Handle back/forward navigation
  window.onpopstate = function (event) {
    if (event.state && event.state.page) {
      self.loadPage(event.state.page);
    }
  };

  // Initial load
  const initial = location.pathname.replace('/', '') || 'home';
  self.loadPage(initial);
})();

