<template>
  <div class="proposals_detail_wrap">
    <div class="proposals_title_wrap">
      <p :class="proposalsDetailWrap" style="margin-bottom:0;">
        <span class="proposals_detail_title">Service</span>
        <span class="proposals_detail_wrap_hash_var">{{`#${name}`}}</span>
      </p>
    </div>
    <div :class="proposalsDetailWrap">
      <p class="proposals_information_content_title">Service Definition Information</p>
      <div class="proposals_detail_information_wrap">
        <div class="information_props_wrap">
          <span class="information_props">Author :</span>
          <span v-show="author !== '--'" class="information_value information_show_trim jump_route" @click="jumpRoute(`/address/1/${author}`)">{{author}}</span>
          <span v-show="author == '--'" class="information_value information_show_trim ">{{author}}</span>
        </div>
        <div class="information_props_wrap">
          <span class="information_props">Chain Id :</span>
          <span class="information_value information_show_trim">
            <pre class="information_pre">{{chainId}}</pre>
          </span>
        </div>
        <div class="information_props_wrap">
          <span class="information_props">Author Description :</span>
          <span class="information_value">{{authorDescription}}</span>
        </div>
        <div class="information_props_wrap">
          <span class="information_props">Description :</span>
          <span class="information_value">{{description}}</span>
        </div>
        <div class="parameter_container">
          <div class="information_props_wrap">
            <span class="information_props">IDL Content :</span>
            <textarea :rows="textareaRows" readonly spellcheck="false" class="parameter_detail_content">{{iDLContent}}
            </textarea>
          </div>
        </div>
      </div>
    </div>
    <div :class="proposalsDetailWrap">
      <p class="proposals_information_content_title" style='border-bottom:none !important;'>Service Bond Detail</p>
    </div>
    <div :class="proposalsDetailWrap">
      <div class="proposals_detail_table_wrap">
        <spin-component :showLoading="showLoading"/>
        <blocks-list-table :items="bondRecord" :type="'ServiceBind'" :showNoData="showNoData" :min-width="tableMinWidth"></blocks-list-table>
        <div v-show="showNoData" class="no_data_show">
          No Data
        </div>
      </div>
    </div>

    <div :class="proposalsDetailWrap">
      <p class="proposals_information_content_title" style='border-bottom:none !important;'>Service Transactions</p>
    </div>
    <div :class="proposalsDetailWrap">
      <div class="proposals_detail_table_wrap">
        <spin-component :showLoading="showLoading"/>
        <blocks-list-table :items="invocationRecord" :type="'ServiceInvocation'" :showNoData="showNoData" :min-width="tableMinWidth"></blocks-list-table>
        <div v-show="showNoData" class="no_data_show">
          No Data
        </div>
      </div>
    </div>

  </div>
</template>

<script>
  import Tools from '../util/Tools';
  import Service from "../util/axios"
  import BlocksListTable from './table/BlocksListTable.vue';
  import SpinComponent from './commonComponents/SpinComponent';
  import Constant from "../constant/Constant"
  export default {
    components: {
      BlocksListTable,
      SpinComponent
    },
    data() {
      return {
        devicesWidth: window.innerWidth,
        proposalsDetailWrap: 'personal_computer_transactions_detail',
        bondRecord: [],
        invocationRecord: [],
        showLoading:false,
        showNoData: false,
        count: "",
        chainId: "",
        name: "",
        description: "",
        author: "",
        authorDescription: "",
        tableMinWidth: "",
        textareaRows: '2',
        iDLContent: '',
        parameterValue: ''
      }
    },
    beforeMount() {
      Tools.scrollToTop();
      if (Tools.currentDeviceIsPersonComputer()) {
        this.proposalsDetailWrap = 'personal_computer_transactions_detail_wrap';
      } else {
        this.proposalsDetailWrap = 'mobile_transactions_detail_wrap';
      }
    },
    mounted() {
      this.getIService();
    },
    methods: {
      getIService() {
        this.showLoading = true;
        let url = `/api/service/${this.$route.params.svcName}/${this.$route.params.defChainId}`;
        Service.http(url).then((data) => {
          this.showLoading = false;
          if(data){
            this.name = data.name;
            this.author = data.author;
            this.chainId = data.chain_id;
            this.authorDescription = data.author_description;
            this.description = data.description;
            this.bondRecord = data.svc_bind_list.map(item =>{
              return {
                'Hash' : item.hash,
                'Bond Chain Id' : item.bind_chain_id,
                'Provider' : item.provider,
                'Bond Type' : item.binding_type,
                'Prices' : item.price,
                'Status' : item.available,
              }
            });

            this.invocationRecord = data.svc_tx.map(item =>{
              return {
                'Hash' : item.hash,
                'Tx Type' : item.tx_type,
                'Send Address': item.send_addr,
                'Receive Address': item.receive_addr,
                'Height': item.height,
                'Time': item.time,
              }
            });
          }else {
              this.showNoData = false;
              this.bondRecord = [{
                'Hash' : '',
                'Bond Chain Id' : '',
                'Provider' : '',
                'Bond Type' : '',
                'Prices' : '',
                'Status' : '',
              }];
            this.invocationRecord = [{
              'Hash': "",
              'Tx Type': "",
              'Send Address': "",
              'Receive Address': "",
              'Height': "",
              'Time': "",
            }];
              this.showNoData = true
            }
        }).catch(e => {
          this.bondRecord = [{
            'Hash' : '',
            'Bond Chain Id' : '',
            'Provider' : '',
            'Bond Type' : '',
            'Prices' : '',
            'Status' : '',
          }];
          this.invocationRecord = [{
            'Hash': "",
            'Tx Type': "",
            'Send Address': "",
            'Receive Address': "",
            'Height': "",
            'Time': "",
          }];
          this.showLoading = false;
          this.showNoData = true
        })

      },
      jumpRoute(path) {
        this.$router.push(path);
      }

    }
  }
</script>

<style scoped lang="scss">
  @import '../style/mixin.scss';

  .proposals_detail_wrap {
    @include flex;
    @include pcContainer;
      font-size: 0.14rem;
    .proposals_title_wrap {
      width: 100%;
      //border-bottom: 1px solid #d6d9e0;
    @include flex;
    @include pcContainer;
    .personal_computer_transactions_detail_wrap {
      @include flex;
    }
  .mobile_transactions_detail_wrap {
    @include flex;
    flex-direction: column;
    .proposals_detail_information_wrap{
      .parameter_container{
        .information_props_wrap{
          .parameter_detail_content{
            width: 90%;
            margin-right:20%;
            background: #EEE;
          }
        }
      }
    }
    }
  }
  .personal_computer_transactions_detail_wrap {
    width: 100%!important;
    .proposals_information_content_title {
      padding-left: 0.2rem !important;
      height: 0.5rem !important;
      line-height: 0.5rem !important;
      font-size: 0.18rem !important;
      color: #000000;
      margin-bottom: 0;
      @include fontWeight;
      border-bottom:1px solid #d6d9e0 !important;
    }
    @include pcCenter;
      .proposals_detail_information_wrap {
        margin-top: 0.21rem;
        margin-left: 0.2rem;
      .information_props_wrap {
        @include flex;
        margin-bottom:0.08rem;
    .information_props {
      min-width: 1.5rem;
    }
    .flag_item_left {
      display: inline-block;
      width: 0.2rem;
      height: 0.17rem;
      background: url('../assets/left.png') no-repeat 0 1px;
      margin-right: 0.05rem;
      cursor: pointer;
    }
    .flag_item_left_disabled {
      display: inline-block;
      width: 0.2rem;
      height: 0.17rem;
      margin-right: 0.05rem;
      cursor: pointer;
      background: url('../assets/left_disabled.png') no-repeat 0 1px;
    }
    .flag_item_right {
      display: inline-block;
      width: 0.2rem;
      height: 0.17rem;
      background: url('../assets/right.png') no-repeat 0 0;
      margin-left: 0.05rem;
      cursor: pointer;
    }
    .flag_item_right_disabled {
      display: inline-block;
      width: 0.2rem;
      height: 0.17rem;
      background: url('../assets/right_disabled.png') no-repeat 0 0;
      margin-left: 0.05rem;
      cursor: pointer;
      }
    }
  }
  .proposals_detail_table_wrap {
    margin-bottom: 0.2rem;
    width: 100%;
    overflow-x: auto;
    .no_data_show {
      @include flex;
        justify-content: center;
        border-top: 0.01rem solid #eee;
        border-bottom: 0.01rem solid #eee;
        font-size: 0.14rem;
        height: 0.5rem;
        align-items: center;
      }
    }

    .proposals_detail_title {
      height: 0.61rem;
      line-height: 0.61rem;
      font-size: 0.22rem;
      color: #000000;
      margin-right: 0.2rem;
      @include fontWeight;
      margin-left: 0.2rem;
    }
    .proposals_detail_wrap_hash_var {
      height: 0.61rem;
      line-height: 0.61rem;
      font-size: 0.22rem;
      color: #a2a2ae;
    }
  }

  .mobile_transactions_detail_wrap {
    width: 100%;
    @include flex;
    flex-direction: column;
    padding-left: 0.1rem;
    .proposals_detail_wrap_hash_var{
      color: #a2a2ae;
    }
    .proposals_information_content_title {
      height: 0.5rem !important;
      line-height: 0.5rem !important;
      font-size: 0.18rem !important;
      color: #000000;
      margin-bottom: 0;
      @include fontWeight;
      border-bottom: 1px solid #d6d9e0 !important;
    }
    .proposals_detail_table_wrap {
      width: 100%;
      overflow-x: auto;
      -webkit-overflow-scrolling:touch;
      margin-bottom:0.4rem;
    .no_data_show {
      @include flex;
        justify-content: center;
        border-top: 0.01rem solid #eee;
        border-bottom: 0.01rem solid #eee;
        font-size: 0.14rem;
        height: 3rem;
        align-items: center;
      }
    }
    .proposals_detail_information_wrap {

      .information_props_wrap {
        @include flex;
          flex-direction: column;
          margin-bottom: 0.05rem;
        .information_value {
          overflow-x: auto;
          -webkit-overflow-scrolling:touch;
        }
        .flag_item_left {
          display: inline-block;
          width: 0.2rem;
          height: 0.17rem;
          background: url('../assets/left.png') no-repeat 0 1px;
          margin-right: 0.05rem;
          cursor: pointer;
        }
        .flag_item_left_disabled {
          display: inline-block;
          width: 0.2rem;
          height: 0.17rem;
          margin-right: 0.05rem;
          cursor: pointer;
          background: url('../assets/left_disabled.png') no-repeat 0 1px;
        }
        .flag_item_right {
          display: inline-block;
          width: 0.2rem;
          height: 0.17rem;
          background: url('../assets/right.png') no-repeat 0 0;
          margin-left: 0.05rem;
          cursor: pointer;
        }
        .flag_item_right_disabled {
          display: inline-block;
          width: 0.2rem;
          height: 0.17rem;
          background: url('../assets/right_disabled.png') no-repeat 0 0;
          margin-left: 0.05rem;
          cursor: pointer;
        }
      }
    }
      .proposals_detail_title {
        height: 0.3rem;
        line-height: 0.3rem;
        font-size: 0.22rem;
        color: #000000;
        margin-right: 0.02rem;
        @include fontWeight;
      }
      .transactions_detail_wrap_hash_var {
        overflow-x: auto;
        -webkit-overflow-scrolling:touch;
        height: 0.3rem;
        line-height: 0.3rem;
        font-size: 0.22rem;
        color: #a2a2ae;
      }
      .vote-details-content{
        width: 100%;
        overflow-x: auto;
        border-top: 1px solid #d6d9e0;
        display: flex;
        justify-content: space-between;
        height: 0.62rem;
        line-height: 0.62rem;
        .vote_content_container{
          min-width: 150%;
          display: flex;
          justify-content: space-between;
        }
      }
    }
  }
  .information_show_trim{
    color: #a2a2ae;
  }
  .information_value{
    color: #a2a2ae;
    word-break: break-all;
  }
  .vote-details-content{
    border-top: 1px solid #d6d9e0;
    display: flex;
    justify-content: space-between;
    height: 0.62rem;
    line-height: 0.62rem;
  }
  .total_num{
    font-size: 0.14rem;
    color: #a2a2ae;
    margin-left: 0.2rem;
  }
  .voting_options{
    display: flex;
    color: #a2a2ae;
    span{
      font-size: 0.14rem;
      color: #000;
      @include fontWeight;
      padding: 0 0.18rem;
    }
  }
  .information_show_trim{
    white-space: pre-wrap ;
  }
  .jump_route {
    color: #3598db;
    cursor: pointer;
  }
  .vote_content_container{
    min-width: 100%;
    display: flex;
    justify-content: space-between;
  }
  pre{
    font-family: Arial !important;
  }
  .information_link{
    color: #3498db !important;
  }
  .parameter_detail_content{
    box-sizing: border-box;
    padding: 0.1rem;
    width: 97%;
    margin-right:20%;
    background: #EEE !important;
  }
</style>
