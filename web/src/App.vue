<template>
  <div id="app">
    <b-navbar :fixed-top="true">
        <template slot="brand">
          <div class="logo">
            <router-link :to="{name: 'dashboard'}">
              <img class="full" src="@/assets/logo.png"/>
              <img class="favicon" src="@/assets/logo.png"/>
            </router-link>
          </div>
        </template>
        <template slot="start">
            <b-navbar-item tag="router-link" :to="{ name : 'dashboard'}"
              label="Dashboard">
                Dashboard
            </b-navbar-item>
        </template>
         <template slot="start">
            <b-navbar-item tag="router-link" :to="{ name : 'settings'}"
              label="Settings">
                Settings
            </b-navbar-item>
        </template>
        <template slot="end">
            <b-navbar-item href="https://github.com/geekwhocodes/innocent-act">
                Github
            </b-navbar-item>
        </template>
    </b-navbar>

    <div class="wrapper">
      <div class="main">
          <router-view :key="$route.fullPath" />
      </div>
    </div>
    <b-loading v-if="!isLoaded" active>
        <div class="has-text-centered">
          <h1 class="title">Oops</h1>
          <p>
            Can't connect to the backend.<br />
            Make sure the server is running and refresh this page.
          </p>
        </div>
    </b-loading>
  </div>
</template>

<script>
import Vue from 'vue';

export default Vue.extend({
  name: 'App',

  data() {
    return {
      isLoaded: true,
    };
  },
  methods: {
  },

  computed: {
    version() {
      return process.env.VUE_APP_VERSION;
    },
  },

  mounted() {
    // Lists is required across different views. On app load, fetch the lists
    // and have them in the store.
    setTimeout(() => {
      this.isLoaded = true;
    }, 500);
  },
});
</script>

<style lang="scss">
  @import "assets/style.scss";
  // @import "assets/icons/unicons.css";
</style>
