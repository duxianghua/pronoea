<template>
  <div class="app-container">
    <el-row style="margin-bottom: 20px;">
      <el-button @click="AddContactGroup">
        Add ContactGroup
      </el-button>
    </el-row>
    <!-- Table -->
    <el-table
      v-loading="listLoading"
      :data="items"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
    >
      <el-table-column label="Name" width="260">
        <template slot-scope="scope">
          {{ scope.row.metadata.name }}
        </template>
      </el-table-column>
      <el-table-column label="Projects" align="center">
        <template slot-scope="scope">
          {{ scope.row.spec.projects }}
        </template>
      </el-table-column>
      <el-table-column label="Members" align="center">
        <template slot-scope="scope">
          <!-- <p v-for="v in scope.row.spec.members">
            {{ v }}
          </p> -->
          <div style="vertical-align: sub;">
            <el-tag
              v-for="member,key in scope.row.spec.members"
              :key="key"
              size="medium"
              style="float: left;margin: 5px 5px 5px 5px;"
              effect="plain"
            >
              {{ member }}
            </el-tag>
          </div>
        </template>
      </el-table-column>

      <el-table-column align="center" prop="created_at" label="CreateTime" width="200">
        <template slot-scope="scope">
          <!-- <i class="el-icon-time"/> -->
          <span>{{ scope.row.metadata.creationTimestamp }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="Aciton" width="200">
        <template slot-scope="scope">
          <!-- <el-button type="primary" @click="probeStatus(scope.row.metadata)">Check</el-button> -->
          <el-button type="text" icon="el-icon-edit-outline" circle size="mini" @click="editContactGroup(scope.row)">EDIT</el-button>
          <el-button type="text" icon="el-icon-delete" circle size="mini" @click="delContactGroup(scope.row)">DELETE</el-button>
        </template>

      </el-table-column>
    </el-table>
    <ContactGroupFrom :is-active.sync="formShow" :form-data="formItems" :is-edit="formEdit" :call-back="fetchData" />
  </div>
</template>

<script>
import { ListContactGroup, DeleteContactGroup } from '@/api/ContactGroup'
import ContactGroupFrom from './from'

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
    ContactGroupFrom
  },
  data() {
    return {
      centerDialogVisible: false,
      statusContext: {},
      items: null,
      listLoading: true,
      showDrawer: false,
      formItems: {},
      formShow: false,
      formEdit: false,
      contactGroupItem: {
        'kind': 'ContactGroup',
        'apiVersion': 'pronoea.io/v1',
        'metadata': {
          'name': '',
          'labels': {}
        },
        'spec': {
          'members': [''],
          'projects': ''
        }
      }
    }
  },
  watch: {
    showDrawer: {
      handler: function(val, oldval) {
        console.log('showDrawer: ' + val)
      },
      immediate: true
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      ListContactGroup().then(response => {
        this.items = response.items
        this.listLoading = false
      }).catch(err => {
        console.log(err)
      })
    },
    delContactGroup(row) {
      DeleteContactGroup(row.metadata.name).then(res => {
        this.items.splice(row.$index, 1)
      }).catch(err => {
        console.log(err)
      })
    },
    editContactGroup(row) {
      this.formItems = row
      this.formShow = true
      this.formEdit = true
    },
    AddContactGroup() {
      this.formItems = this.contactGroupItem
      this.formShow = true
      this.formEdit = false
    }
  }
}
</script>
