<template>
  <section>
  <div>Decent Dashboard</div>
  </section>
</template>

<script>
// @ is an alias to /src

export default {
  name: 'Dashboard',
  components: {
  },

  data() {
    return {
      isListLoading: true,
    };
  },

  mounted() {
    setTimeout(() => {
      this.isListLoading = false;
    }, 2000);
  },
};
</script>
