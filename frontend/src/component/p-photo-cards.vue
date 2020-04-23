<template>
    <v-container grid-list-xs fluid class="pa-2 p-photos p-photo-cards">
        <v-card v-if="photos.length === 0" class="p-photos-empty secondary-light lighten-1 ma-1" flat>
            <v-card-title primary-title>
                <div>
                    <h3 class="title mb-3">
                        <translate>No photos matched your search</translate>
                    </h3>
                    <div>
                        <translate>Try using other terms and search options such as category, country and camera.
                        </translate>
                    </div>
                </div>
            </v-card-title>
        </v-card>
        <v-layout row wrap class="p-results">
            <v-flex
                    v-for="(photo, index) in photos"
                    :key="index"
                    class="p-photo"
                    xs12 sm6 md4 lg3 d-flex
                    v-bind:class="{ 'is-selected': $clipboard.has(photo) }"
            >
                <v-hover>
                    <v-card tile slot-scope="{ hover }"
                            @contextmenu="contextMenu($event, photo, index)"
                            :dark="$clipboard.has(photo)"
                            :class="$clipboard.has(photo) ? 'elevation-10 ma-0 accent darken-1 white--text' : 'elevation-0 ma-1 accent lighten-3'">
                        <v-img
                                :src="photo.getThumbnailUrl('tile_500')"
                                aspect-ratio="1"
                                v-bind:class="{ selected: $clipboard.has(photo) }"
                                style="cursor: pointer;"
                                class="accent lighten-2"
                                v-longclick="longClick"
                                @click="onClick($event, photo, index)"
                        >
                            <v-layout
                                    slot="placeholder"
                                    fill-height
                                    align-center
                                    justify-center
                                    ma-0
                            >
                                <v-progress-circular indeterminate color="accent lighten-5"></v-progress-circular>
                            </v-layout>

                            <v-btn v-if="hover || selection.length > 0" :flat="!hover" :ripple="false"
                                   icon large absolute
                                   :class="$clipboard.has(photo) ? 'p-photo-select' : 'p-photo-select opacity-50'"
                                   @click.stop.prevent="onSelect($event, photo, index)">
                                <v-icon v-if="selection.length && $clipboard.has(photo)" color="white"
                                        class="t-select t-on">check_circle
                                </v-icon>
                                <v-icon v-else color="accent lighten-3" class="t-select t-off">radio_button_off</v-icon>
                            </v-btn>

                            <v-btn :flat="!hover" :ripple="false"
                                   icon large absolute
                                   :class="photo.PhotoFavorite ? 'p-photo-like opacity-75' : 'p-photo-like opacity-50'"
                                   @click.stop.prevent="photo.toggleLike()">
                                <v-icon v-if="photo.PhotoFavorite" color="white" class="t-like t-on">favorite</v-icon>
                                <v-icon v-else color="accent lighten-3" class="t-like t-off">favorite_border</v-icon>
                            </v-btn>

                            <v-btn v-if="photo.Files.length > 1" :flat="!hover" :ripple="false"
                                   icon large absolute class="p-photo-merged opacity-75"
                                   @click.stop.prevent="openPhoto(index, true)">
                                <v-icon color="white" class="action-burst">burst_mode</v-icon>
                            </v-btn>
                        </v-img>

                        <v-card-title primary-title class="pa-3 p-photo-desc" style="user-select: none;">
                            <div>
                                <h3 class="body-2 mb-2" :title="photo.PhotoTitle">
                                    <button @click.exact="editPhoto(index)">
                                        {{ photo.PhotoTitle | truncate(80) }}
                                    </button>
                                </h3>
                                <div class="caption">
                                    <button @click.exact="editPhoto(index)">
                                        <v-icon size="14">date_range</v-icon>
                                        {{ photo.getDateString() }}
                                    </button>
                                    <br/>
                                    <button @click.exact="editPhoto(index)">
                                        <v-icon size="14">photo_camera</v-icon>
                                        {{ photo.getCamera() }}
                                    </button>
                                    <br/>
                                    <button @click.exact="openLocation(index)" v-if="showLocation && photo.LocationID">
                                        <v-icon size="14">location_on</v-icon>
                                        {{ photo.getLocation() }}
                                    </button>
                                </div>
                            </div>
                        </v-card-title>
                    </v-card>
                </v-hover>
            </v-flex>
        </v-layout>
    </v-container>
</template>
<script>
    export default {
        name: 'p-photo-cards',
        props: {
            photos: Array,
            selection: Array,
            openPhoto: Function,
            editPhoto: Function,
            openLocation: Function,
            album: Object,
        },
        data() {
            return {
                showLocation: this.$config.settings().features.places,
                wasLong: false,
            };
        },
        methods: {
            longClick() {
                this.wasLong = true;
            },
            onSelect(ev, model, index) {
                if (ev.shiftKey) {
                    this.selectRange(index);
                } else {
                    this.$clipboard.toggle(model);
                }

                this.wasLong = false;
            },
            onClick(ev, model, index) {
                if (this.wasLong || this.selection.length > 0) {
                    ev.preventDefault();
                    ev.stopPropagation();

                    if (this.wasLong || ev.shiftKey) {
                        this.selectRange(index);
                    } else {
                        this.$clipboard.toggle(model);
                    }
                } else {
                    this.openPhoto(index, false);
                }

                this.wasLong = false;
            },
            contextMenu(ev, model, index) {
                if (this.$isMobile) {
                    ev.preventDefault();
                    ev.stopPropagation();

                    if (this.wasLong) {
                        this.selectRange(index);
                    } else {
                        this.$clipboard.toggle(model);
                    }
                }

                this.wasLong = false;
            },
            selectRange(index) {
                this.$clipboard.addRange(index, this.photos);
            }
        }
    };
</script>