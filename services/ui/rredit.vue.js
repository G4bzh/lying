export default {
  name: "RrEdit",
  template: `
    <div>
      <md-card class="md-elevation-9">

        <md-card-actions>
         <md-button class="md-icon-button md-dense md-raised md-accent" @click="doRemove" v-if="rrindex != undefined"><md-icon >clear</md-icon></md-button>
         <md-button class="md-icon-button md-dense md-raised md-primary" @click="doAdd" v-else><md-icon >add</md-icon></md-button>
        </md-card-actions>


        <md-card-content class="md-layout">

           <md-field class="md-layout-item md-size-30">
            <label>Name</label>
            <md-input v-model="rrname_" placeholder="Name"></md-input>
           </md-field>

           <md-field class="md-layout-item md-size-10">
            <label>TTL</label>
            <md-input v-model="rrttl_" placeholder="TTL"></md-input>
           </md-field>

           <md-field class="md-layout-item md-size-10">
            <label>Class</label>
            <md-select v-model="rrclass_" md-dense placeholder="Class">
                <md-option value="IN">IN</md-option>
                <md-option value="CH">CH</md-option>
                <md-option value="HS">HS</md-option>
                <md-option value="NONE">NONE</md-option>
                <md-option value="*">ANY</md-option>
            </md-select>
          </md-field>

          <md-field class="md-layout-item md-size-10">
          <label>Type</label>
           <md-select v-model="rrtype_" md-dense placeholder="Type">
             <md-option value="*">*</md-option>
             <md-option value="A">A</md-option>
             <md-option value="AAAA">AAAA</md-option>
             <md-option value="AFSDB">AFSDB</md-option>
             <md-option value="APL">APL</md-option>
             <md-option value="ATMA">ATMA</md-option>
             <md-option value="AVC">AVC</md-option>
             <md-option value="AXFR">AXFR</md-option>
             <md-option value="CAA">CAA</md-option>
             <md-option value="CDNSKEY">CDNSKEY</md-option>
             <md-option value="CDS">CDS</md-option>
             <md-option value="CERT">CERT</md-option>
             <md-option value="CNAME">CNAME</md-option>
             <md-option value="CSYNC">CSYNC</md-option>
             <md-option value="DHCID">DHCID</md-option>
             <md-option value="DLV">DLV</md-option>
             <md-option value="DNAME">DNAME</md-option>
             <md-option value="DNSKEY">DNSKEY</md-option>
             <md-option value="DOA">DOA</md-option>
             <md-option value="DS">DS</md-option>
             <md-option value="EID">EID</md-option>
             <md-option value="EUI48">EUI48</md-option>
             <md-option value="EUI64">EUI64</md-option>
             <md-option value="GID">GID</md-option>
             <md-option value="GPOS">GPOS</md-option>
             <md-option value="HINFO">HINFO</md-option>
             <md-option value="HIP">HIP</md-option>
             <md-option value="IPSECKEY">IPSECKEY</md-option>
             <md-option value="ISDN">ISDN</md-option>
             <md-option value="IXFR">IXFR</md-option>
             <md-option value="KEY">KEY</md-option>
             <md-option value="KX">KX</md-option>
             <md-option value="L32">L32</md-option>
             <md-option value="L64">L64</md-option>
             <md-option value="LOC">LOC</md-option>
             <md-option value="LP">LP</md-option>
             <md-option value="MAILB">MAILB</md-option>
             <md-option value="MB">MB</md-option>
             <md-option value="MG">MG</md-option>
             <md-option value="MINFO">MINFO</md-option>
             <md-option value="MR">MR</md-option>
             <md-option value="MX">MX</md-option>
             <md-option value="NAPTR">NAPTR</md-option>
             <md-option value="NID">NID</md-option>
             <md-option value="NIMLOC">NIMLOC</md-option>
             <md-option value="NINFO">NINFO</md-option>
             <md-option value="NS">NS</md-option>
             <md-option value="NSAP">NSAP</md-option>
             <md-option value="NSAP-PTR">NSAP-PTR</md-option>
             <md-option value="NSEC">NSEC</md-option>
             <md-option value="NSEC3">NSEC3</md-option>
             <md-option value="NSEC3PARAM">NSEC3PARAM</md-option>
             <md-option value="NULL">NULL</md-option>
             <md-option value="OPENPGPKEY">OPENPGPKEY</md-option>
             <md-option value="OPT">OPT</md-option>
             <md-option value="Private use">Private use</md-option>
             <md-option value="PTR">PTR</md-option>
             <md-option value="PX">PX</md-option>
             <md-option value="Reserved">Reserved</md-option>
             <md-option value="RKEY">RKEY</md-option>
             <md-option value="RP">RP</md-option>
             <md-option value="RRSIG">RRSIG</md-option>
             <md-option value="RT">RT</md-option>
             <md-option value="SIG">SIG</md-option>
             <md-option value="SINK">SINK</md-option>
             <md-option value="SMIMEA">SMIMEA</md-option>
             <md-option value="SOA">SOA</md-option>
             <md-option value="SPF">SPF</md-option>
             <md-option value="SRV">SRV</md-option>
             <md-option value="SSHFP">SSHFP</md-option>
             <md-option value="TA">TA</md-option>
             <md-option value="TALINK">TALINK</md-option>
             <md-option value="TKEY">TKEY</md-option>
             <md-option value="TLSA">TLSA</md-option>
             <md-option value="TSIG">TSIG</md-option>
             <md-option value="TXT">TXT</md-option>
             <md-option value="UID">UID</md-option>
             <md-option value="UINFO">UINFO</md-option>
             <md-option value="UNSPEC">UNSPEC</md-option>
             <md-option value="URI">URI</md-option>
             <md-option value="WKS">WKS</md-option>
             <md-option value="X25">X25</md-option>
             <md-option value="ZONEMD">ZONEMD</md-option>
           </md-select>
         </md-field>

         <md-field class="md-layout-item">
          <label>Data</label>
          <md-input v-model="rrdata_" placeholder="Data"></md-input>
         </md-field>

        </md-card-content>

      </md-card>
    </div>
  `,
  props: {
    rrindex: Number,
    rrname: String,
    rrttl: Number,
    rrclass: String,
    rrtype: String,
    rrdata: String
  },
  data: function() {
    return {
      rrname_: this.rrname,
      rrttl_: this.rrttl,
      rrclass_: this.rrclass,
      rrtype_: this.rrtype,
      rrdata_: this.rrdata
    }
  },
  methods: {
    doRemove: function(event) {
      this.$emit('rr-remove', this.rrindex);
      return;
    },
    doAdd: function(event) {
      this.$emit('rr-add',{
        "name": this.rrname_,
        "rttl": this.rrttl_,
        "class": this.rrclass_,
        "type": this.rrtype_,
        "rdata": this.rrdata_
      });

      this.rrname_ = undefined;
      this.rrttl_ = undefined;
      this.rrclass_ = undefined;
      this.rrtype_ = undefined;
      this.rrdata_ = undefined;

      return;
    }
  }
};
