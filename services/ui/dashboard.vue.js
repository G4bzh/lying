export default {
  name: "Dashboard",
  template: `
  <div>
    <md-content class="md-elevation-3">

      <div class="md-title">Dashboard </div>

      Hello {{ client }}

    </md-content>

  </div>
  `,
  data: function() {
    return {
      client: ""
    }
  }
};
