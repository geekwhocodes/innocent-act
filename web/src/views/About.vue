<template>
  <section class="dashboard content">
    <div class="container" style="position: relative">
      <b-loading v-if="isLoadingUserList" active :is-full-page="false" />
        <b-tabs>
            <b-tab-item label="Table">
                <b-table
                    :data="users"
                    :columns="columns"
                    :checked-rows.sync="checkedRows"
                    checkable
                    :checkbox-position="checkboxPosition">

                    <template slot="bottom-left">
                        <b>Total checked</b>: {{ checkedRows.length }}
                    </template>
                </b-table>
            </b-tab-item>

            <b-tab-item label="Checked rows">
                <pre>{{ checkedRows }}</pre>
            </b-tab-item>
        </b-tabs>
    </div>
  </section>
</template>

<script>
// @ is an alias to /src
import axios from 'axios';

export default {
  name: 'Settings',
  components: {
  },

  data() {
    return {
      users: [],
      errors: [],
      isLoadingUserList: true,
      checkboxPosition: 'left',
      checkedRows: [],
      columns: [
        {
          field: 'ID',
          label: 'ID',
          numeric: true,
        },
        {
          field: 'name',
          label: 'Name',
        },
        {
          field: 'email',
          label: 'Email',
          centered: false,
        },
        {
          field: 'website',
          label: 'Webiste',
          centered: false,
        },
      ],
    };
  },
  async mounted() {
    try {
      axios.defaults.baseURL = 'http://localhost:8083';
      const response = await axios.get('/api/users');
      this.users = response.data;
      console.log(response.data);
      this.isLoadingUserList = false;
    } catch (e) {
      this.errors.push(e);
      this.isLoadingUserList = false;
    }
    // this.isLoadingUserList = false;
  },
  computed: {
  },
};
</script>
