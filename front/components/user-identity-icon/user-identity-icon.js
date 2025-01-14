Component({
  properties: {
    identity: {
      type: String,
      value: '',
    },
    noTip: {
      type: Boolean,
      value: false,
    },
  },

  data: {
    identityIcon: {
      organization: 'success',
    },
    identityName: {
      organization: '组织',
    },
  },
});