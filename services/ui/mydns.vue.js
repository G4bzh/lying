export default {
  name: "MyDns",
  template: `
  <div>
    <RouterView :key="$route.path"></RouterView>
  </div>
  `,
  data: function() {
    return {}
  }
};
