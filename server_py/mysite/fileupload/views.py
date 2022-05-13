from django.shortcuts import render
from django.http import HttpResponseRedirect
from .forms import UploadFileForm
from django.views import generic

# Create your views here.
from .handler import handle_uploaded_file

def upload(request):
    # if request.method == 'POST':
    #     form = UploadFileForm(request.POST, request.FILES);
    #     if form.is_valid():
    #         # Todo: path = 'video' + some string in the content of request.

    #         # example: path_tmp = '/video/student_id/screen.webm'
    #         # example: path_tmp = '/video/student_id/camera.webm'

    #         path_tmp = '_placeholder'
    #         handle_uploaded_file(path_tmp, request.FILES['file'])
    #         return HttpResponseRedirect('/success/index.html')
    # else:
    #     form = UploadFileForm()
    # return render(request, 'upload.html', {'form': form})
    if request.method == 'GET':
        return render(request, 'fileupload/client.html')
    elif request.method == 'POST':
        file_obj = request.FILES.get('file')
        file_name = request.POST.get('filename')
        print(file_name)
        handle_uploaded_file(file_name, file_obj)
        return HttpResponseRedirect('success')

def success(request):
    return render(request, 'success/index.html')
