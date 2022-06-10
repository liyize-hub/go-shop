// 同时发送异步代码的次数
let ajaxTimes = 0;
export const request = (params) => {

  ajaxTimes++,
    // 显示加载中 效果
    wx.showLoading({
      title: "加载中",
      mask: true
    });

  // 定义公共的url
  return new Promise((resolve, reject) => {
    wx.request({
      ...params,
      success: (result) => {
        resolve(result.data);
      },
      fail: (err) => {
        reject(err)
      },
      complete: () => {
        ajaxTimes--;
        if (ajaxTimes === 0) {
          //  关闭正在等待的图标
          wx.hideLoading();
        }
      }
    });

  })
}


/**
 *  promise 形式  showModal
 * @param {object} param0 参数
 */
export const showModal = ({ content }) => {
  return new Promise((resolve, reject) => {
    wx.showModal({
      title: '提示',
      content: content,
      success: (res) => {
        resolve(res);
      },
      fail: (err) => {
        reject(err);
      }
    })
  })
}

/**
 *  promise 形式  showToast
 * @param {object} param0 参数
 */
export const showToast = ({ title }) => {
  return new Promise((resolve, reject) => {
    wx.showToast({
      title: title,
      icon: 'none',
      success: (res) => {
        resolve(res);
      },
      fail: (err) => {
        reject(err);
      }
    })
  })
}

/**
 *  promise 形式  getUserProfile
 * @param {object} param0 参数
 */
 export const getUserProfile = ({ desc }) => {
  return new Promise((resolve, reject) => {
    wx.getUserProfile({
      desc: desc,
      success: (res) => {
        resolve(res);
      },
      fail: (err) => {
        reject(err);
      }
    })
  })
}