
const state = {
  showDrawer: false
}

const mutations = {
  CHANGE_showSettings: (state, { key, value }) => {
    if (state.hasOwnProperty(key)) {
      state[key] = value
    }
  },
  ChangeShowDrawer() {
    state.showDrawer = !state.showDrawer
  }
}

export default {
  state,
  mutations
}
