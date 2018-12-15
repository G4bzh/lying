export default {
  name: "MyDnsZone",
  template: `
  <div>

      Zone {{ $route.params.zone }}

  </div>
  `,
  data: function() {
    return {}
  }
};
