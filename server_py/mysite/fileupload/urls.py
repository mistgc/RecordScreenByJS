from django.urls import path

from . import views

urlpatterns = [
    path('', views.upload),
    path('file', views.upload),
    path('success', views.success),
]
