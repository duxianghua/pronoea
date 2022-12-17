<template>
  <div>
    <el-row v-for="(item,index) in labels" :key="index">
      <el-col :span="9">
        <el-input v-model="labels[index].key" :disabled="item.key === 'project'" @input="changeEvent" />
      </el-col>
      <el-col :span="9">
        <el-input v-model="labels[index].value" :disabled="item.key === 'project'" style="margin-left: 10px;" @input="changeEvent" />
      </el-col>
      <el-col :span="6">
        <el-button-group>
          <el-button v-show="labels.length - 1 == index" type="text" icon="el-icon-document-add" :disabled="addButtonState" circle size="medium" @click="addLabel(labels)">Add</el-button>
          <el-button type="text" icon="el-icon-delete" circle size="medium" :disabled="item.key === 'project'" @click="delLabel(index, labels)">Delete</el-button>
        </el-button-group>
        <!-- <svg @click="addLabel(labels)" viewBox="0 0 24 24" width="32" height="32" v-show="labels.length - 1 == index">
          <path fill="none" d="M0 0h24v24H0z" />
            <path
              fill="green"
              d="M11 11V7h2v4h4v2h-4v4h-2v-4H7v-2h4zm1 11C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm0-2a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"
            />
        </svg>
        <svg @click="delLabel(index, labels)" viewBox="0 0 24 24" width="32" height="32">
          <path fill="none" d="M0 0h24v24H0z" />
            <path
              fill="red"
              d="M12 22C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm0-2a8 8 0 1 0 0-16 8 8 0 0 0 0 16zm0-9.414l2.828-2.829 1.415 1.415L13.414 12l2.829 2.828-1.415 1.415L12 13.414l-2.828 2.829-1.415-1.415L10.586 12 7.757 9.172l1.415-1.415L12 10.586z"
            />
        </svg> -->
      </el-col>
    </el-row>
  </div>

</template>

<script>

export default {
  name: 'Lables',
  model: {
    prop: 'items',
    event: 'change'
  },
  props: {
    items: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      disableAddButton: true
    }
  },
  computed: {
    labels: {
      get() {
        var _labels = []
        for (var key in this.items) {
          if (key === 'project') {
            _labels = [{ key: key, value: this.items[key] }].concat(_labels)
            continue
          }
          _labels.push({ key: key, value: this.items[key] })
        }
        return _labels
      }
    },
    addButtonState: {
      get() {
        return (this.labels[this.labels.length - 1].key === '' || this.labels[this.labels.length - 1].value === '')
      }
    }
  },
  methods: {
    changeEvent(event) {
      var d = {}
      this.labels.forEach(element => {
        d[element['key']] = element['value']
      })
      this.$emit('change', d)
    },
    addLabel(items) {
      if (items[items.length - 1].key !== '') {
        items.push({ key: '', value: '' })
        this.changeEvent()
      }
    },
    delLabel(index, items) {
      if (items.length === 1 && items[0].key !== '' && items[0].value !== '') {
        items.splice(index, 1)
        items.push({ key: '', value: '' })
        this.changeEvent()
      }
      if (items.length > 1) {
        items.splice(index, 1)
        this.changeEvent()
      }
    }
  }

}
</script>

<style scoped>
.svg-icon {
  width: 1em;
  height: 1em;
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
}

.svg-external-icon {
  background-color: currentColor;
  mask-size: cover!important;
  display: inline-block;
}
.el-col{
  margin-bottom: 10px;
  text-align: right;
}
.svg-icon{
  text-align: right;
}
</style>
