import { mount } from "@vue/test-utils"
import '@testing-library/jest-dom'
import Home from "@/views/Home.vue";
import Footer from "@/components/Footer.vue"
import { rootStoreFactory, IRootStoreState } from "@/store/root"
import { RecommendationRaw } from "@/store/recommendation_raw";
import { RecommendationExtra } from "@/store/recommendation_extra";
import { freshSampleRawRecommendation } from "./sample_recommendation";
import vuetify from '@/plugins/vuetify';
import { coreTableStoreStateFactory, ICoreTableStoreState } from '@/store/core_table';
import { IRecommendationsStoreState, recommendationsStoreStateFactory } from '../../src/store/recommendations';
import { Store } from 'vuex';
import CoreTable from '../../src/components/CoreTable.vue';

describe("Footer", () => {
  let tableState: ICoreTableStoreState;
  let recommendation: RecommendationExtra;
  let fakeStore: Store<IRootStoreState>;

  beforeEach(() => {
    tableState = coreTableStoreStateFactory();
    recommendation = new RecommendationExtra(freshSampleRawRecommendation());
    fakeStore = rootStoreFactory();
    fakeStore.commit("recommendationsStore/addRecommendation", recommendation);
  });
  
  it("When there are no recommendations, there is no footer", () => {
    const fakeStore = rootStoreFactory();
    const wrapper = mount(Home, { store: fakeStore, vuetify: vuetify })

    expect(wrapper.find("[data-name=footer]").exists()).toBe(false);
  }),

    it("When there are selected recommendations, there is a footer", () => {
      const fakeStore = rootStoreFactory();
      const wrapper = mount(Home, { store: fakeStore, vuetify: vuetify });
      const selected = new Array<RecommendationExtra>();
      selected.push(recommendation);
      fakeStore.commit("coreTableStore/setSelected", selected);
      console.log(wrapper.findComponent(CoreTable).find(".v-data-table__wrapper > table:nth-child(1) > tbody:nth-child(3) > tr:nth-child(2) > td:nth-child(1) > div:nth-child(1)"));
     // console.log(wrapper.findComponent(Footer).findComponent({ref: "test"}).element);
    expect(wrapper.find(".footer").exists()).toBe(true);
  })
})
