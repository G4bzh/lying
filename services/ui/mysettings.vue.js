export default {
  name: "MySettingss",
  template: `
  <div>
      {{username}} settings
  </div>
  `,
  data: function() {
    return {}
  },
  computed: {
    username: function () {
      const user = JSON.parse(localStorage.getItem('user'));
      if (user) {
        if (user.name) {
          return user.name;
        }
      }
      return "Unknown";
    }
  }
};
