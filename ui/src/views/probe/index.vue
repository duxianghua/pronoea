<template>
  <div class="app-container">
    <div style="margin-bottom: 5px; box-shadow: 0px 0px 10px #dedede; height: 60px; padding:10px">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-input
            v-model="search"
            placeholder="Search"
            prefix-icon="el-icon-search"
          />
        </el-col>
        <el-col :span="8"><div class="grid-content bg-purple" /></el-col>
        <el-col :span="8" style="float:right">
          <el-button type="primary" style="float:right" @click="openFrom({})">
            Add Probe
          </el-button>
        </el-col>
      </el-row>
    </div>
    <!-- Table -->
    <el-table
      :key="listKey"
      v-loading="listLoading"
      :data="list.filter(data => !search || data.metadata.name.toLowerCase().includes(search.toLowerCase()))"
      element-loading-text="Loading"
      border
      fit
      size="nini"
      highlight-current-row
      style="min-height: 540px; box-shadow: 0px 0px 10px #dedede;"
    >
      >
      <el-table-column align="center" label="ID" width="95">
        <template slot-scope="scope">
          {{ scope.$index }}
        </template>
      </el-table-column>
      <el-table-column label="Name" width="260">
        <template slot-scope="scope">
          {{ scope.row.metadata.name }}
        </template>
      </el-table-column>
      <el-table-column label="Targets" align="center">
        <template slot-scope="scope">
          <div v-for="target,key in scope.row.spec.targets" :key="key">
            <el-link type="primary" :href="v" target="_blank">{{ target }}</el-link>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="ProbeType" width="110" align="center">
        <template slot-scope="scope">
          {{ scope.row.spec.module.prober }}
        </template>
      </el-table-column>
      <el-table-column label="Monitoring" width="110" align="center">
        <template slot-scope="scope">
          <el-switch
            v-model="scope.row.spec.pause"
            active-color="#13ce66"
            inactive-color="#ff4949"
            @change="pauseSwitch(scope)"
          />
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="CreateTime" width="200">
        <template slot-scope="scope">
          <!-- <i class="el-icon-time"/> -->
          <span>{{ scope.row.metadata.creationTimestamp }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="Aciton" width="210">
        <template slot-scope="scope">
          <el-button type="primary" size="mini" @click="openFrom(scope.row)">Edit</el-button>
          <el-button type="primary" size="mini" @click="probeStatus(scope.row)">Check</el-button>
          <el-button type="danger" icon="el-icon-delete" size="mini" @click="delProbe(scope)" />
        </template>

      </el-table-column>
    </el-table>
    <probeForm :is-active.sync="showDrawer" :form-data="formData" :call-back="fetchData" />
    <el-dialog
      title=""
      :visible.sync="centerDialogVisible"
      width="80%"
      center
    >
      <template>
        <el-tabs>
          <el-tab-pane v-for="valuea, key, index in statusContext" :key="index" :label="key">
            <el-input
              v-loading="loading"
              type="textarea"
              :rows="18"
              :value="valuea"
            />
          </el-tab-pane>
        </el-tabs>
      </template>
    <!-- <span slot="footer" class="dialog-footer">
      <el-button @click="centerDialogVisible = false">Close</el-button>
      <el-button type="primary" @click="centerDialogVisible = false">确 定</el-button>
    </span> -->
    </el-dialog>
  </div>
</template>

<script>
import { getList, deleteProbe, statusProbe, UpdateProbe } from '@/api/probe'
import probeForm from './form'

export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'gray',
        deleted: 'danger'
      }
      return statusMap[status]
    }
  },
  components: {
    probeForm
  },
  data() {
    return {
      formEdit: false,
      formData: {},
      showDrawer: false,
      centerDialogVisible: false,
      statusContext: {},
      list: [],
      listKey: Math.random(),
      listLoading: true,
      loading: true,
      search: ''
    }
  },
  computed: {
    // showDrawer(){
    //     return this.$store.state.probe.showDrawer
    // }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      getList().then(response => {
        this.list = response.items
        this.listLoading = false
      }).catch(err => {
        console.log(err)
      })
    },
    delProbe(scope) {
      this.$confirm('Are you sure to delete permanently', {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        deleteProbe(scope.row.metadata.name, { namespace: scope.row.metadata.namespace }).then(response => {
          this.list.splice(scope.$index, 1)
          this.$message({
            type: 'success',
            message: '删除成功!'
          })
        }).catch(err => {
          console.log(err)
        })
      })
    },
    probeStatus(row) {
      this.statusContext = {}
      this.centerDialogVisible = true
      this.loading = true
      row.spec.targets.forEach(element => {
        this.statusContext[element] = ''
      })
      this.loading = true
      statusProbe(row.metadata.name, { namespace: row.metadata.namespace }).then(response => {
        this.statusContext = response
        this.loading = false
      }).catch(err => {
        console.log(err)
        this.loading = false
      })
    },
    DrawerSwitch() {
      this.$store.commit('ChangeShowDrawer')
    },
    closeHander() {
      this.DrawerSwitch()
      this.fetchData()
    },
    pauseSwitch(scope) {
      UpdateProbe(scope.row).then(response => {
        this.$set(this.list, scope.$index, response)
      }).catch(err => {
        this.$message(err.$message)
      })
    },
    openFrom(data) {
      // this.formItems = this.contactGroupItem
      this.formData = data
      this.showDrawer = true
      // this.formEdit = false
    }
  }
}
</script>
